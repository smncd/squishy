package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/resources"
	"gitlab.com/smncd/squishy/internal/router"
)

type SharedContext struct {
	s             *filesystem.SquishyFile
	errorTemplate *template.Template
}

func New(s *filesystem.SquishyFile) *http.Server {
	ctx := SharedContext{
		s:             s,
		errorTemplate: template.Must(template.ParseFS(resources.TemplateFS, "templates/error.html")),
	}

	router := router.New(ctx)

	staticFS, err := fs.Sub(resources.StaticFS, "static")
	if err != nil {
		log.Fatalln("error")
	}

	router.NoRoute(notFoundHandler)

	router.StaticFS("/static/", staticFS)

	router.GET("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router,
	}

	return server
}

func notFoundHandler(w http.ResponseWriter, req *http.Request, ctx SharedContext) {
	w.WriteHeader(http.StatusNotFound)
	ctx.errorTemplate.Execute(w, resources.ErrorTemplateData{
		Title:       "Not Found",
		Description: "The link you've accessed does not exist",
	})
}

func handler(w http.ResponseWriter, r *http.Request, ctx SharedContext) {
	path := r.URL.Path

	tmpl := template.Must(template.ParseFS(resources.TemplateFS, "templates/error.html"))

	err := ctx.s.RefetchRoutes()
	if err != nil {
		data := resources.ErrorTemplateData{
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
