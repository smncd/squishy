package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/router"
	"gitlab.com/smncd/squishy/internal/templates"
)

//go:embed static/*
var staticFS embed.FS

type SharedContext struct {
	s             *filesystem.SquishyFile
	errorTemplate *template.Template
}

func New(s *filesystem.SquishyFile) *http.Server {
	ctx := SharedContext{
		s:             s,
		errorTemplate: template.Must(template.ParseFS(templates.FS, "error.html")),
	}

	router := router.New(ctx)

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalln("error")
	}

	router.NoRoute(notFoundHandler)

	router.StaticFS("/static/", static)

	router.GET("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router,
	}

	return server
}

func notFoundHandler(w http.ResponseWriter, req *http.Request, c SharedContext) {
	w.WriteHeader(http.StatusNotFound)
	c.errorTemplate.Execute(w, templates.ErrorPageData{
		Title:       "Not Found",
		Description: "The link you've accessed does not exist",
	})
}

func handler(w http.ResponseWriter, r *http.Request, c SharedContext) {
	path := r.URL.Path

	tmpl := template.Must(template.ParseFS(templates.FS, "error.html"))

	err := c.s.RefetchRoutes()
	if err != nil {
		data := templates.ErrorPageData{
			Title:       "Welp, that's not good",
			Description: "There's been an error on our end, please check back later",
		}
		if c.s.Config.Debug {
			data.Error = err.Error()
		}

		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, data)
		return
	}

	reply, ok := c.s.LookupRoutePath(path)
	if !ok {
		notFoundHandler(w, r, c)
		return
	}

	http.Redirect(w, r, reply, http.StatusPermanentRedirect)
}
