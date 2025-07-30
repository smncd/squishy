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

type RouterContext struct {
	s             *filesystem.SquishyFile
	errorTemplate *template.Template
}

func New(s *filesystem.SquishyFile) *http.Server {
	rCtx := RouterContext{
		s:             s,
		errorTemplate: template.Must(template.ParseFS(templates.FS, "error.html")),
	}

	router := router.New(rCtx)

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalln("error")
	}

	router.Fallback(notFoundHandler)

	router.StaticFS("/static/", static)

	router.GET("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router,
	}

	return server
}
func notFoundHandler(w http.ResponseWriter, req *http.Request, ctx RouterContext) {
	w.WriteHeader(http.StatusNotFound)
	ctx.errorTemplate.Execute(w, templates.ErrorPageData{
		Title:       "Not Found",
		Description: "The link you've accessed does not exist",
	})
}

func handler(w http.ResponseWriter, r *http.Request, ctx RouterContext) {
	path := r.URL.Path

	tmpl := template.Must(template.ParseFS(templates.FS, "error.html"))

	err := ctx.s.RefetchRoutes()
	if err != nil {
		data := templates.ErrorPageData{
			Title:       "Welp, that's not good",
			Description: "There's been an error on our end, please check back later",
		}
		if ctx.s.Config.Debug {
			data.Error = err.Error()
		}

		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, data)
		return
	}

	reply, ok := ctx.s.LookupRoutePath(path)
	if !ok {
		notFoundHandler(w, r, ctx)
		return
	}

	http.Redirect(w, r, reply, http.StatusPermanentRedirect)
}
