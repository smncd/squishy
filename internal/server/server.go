package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/server/router"
	"gitlab.com/smncd/squishy/internal/templates"
)

//go:embed static/*
var staticFS embed.FS

func New(s *filesystem.SquishyFile) *http.Server {
	router, err := router.New(router.RouterContext{S: s})
	if err != nil {
		log.Fatalln("error")
	}

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalln("error")
	}

	router.StaticFS("/static/", static)

	router.GET("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router,
	}

	return server
}

func handler(w http.ResponseWriter, r *http.Request, ctx router.RouterContext) {
	path := r.URL.Path

	tmpl := template.Must(template.ParseFS(templates.FS, "error.html"))

	err := ctx.S.RefetchRoutes()
	if err != nil {
		data := templates.ErrorPageData{
			Title:       "Welp, that's not good",
			Description: "There's been an error on our end, please check back later",
		}
		if ctx.S.Config.Debug {
			data.Error = err.Error()
		}

		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Execute(w, data)
		return
	}

	reply, ok := ctx.S.LookupRoutePath(path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, templates.ErrorPageData{
			Title:       "Not Found",
			Description: "The link you've accessed does not exist",
		})

		return
	}

	http.Redirect(w, r, reply, http.StatusPermanentRedirect)
}
