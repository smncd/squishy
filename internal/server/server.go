package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/logging"
	"gitlab.com/smncd/squishy/internal/resources"
	"gitlab.com/smncd/squishy/internal/router"
)

type SharedContext struct {
	s             *filesystem.SquishyFile
	errorTemplate *template.Template
	logger        *log.Logger
}

func New(s *filesystem.SquishyFile, logger *log.Logger) *http.Server {
	sc := SharedContext{
		s:             s,
		errorTemplate: template.Must(template.ParseFS(resources.TemplateFS, "templates/error.html")),
		logger:        logger,
	}

	router := router.New(sc, logger)

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

func notFoundHandler(w http.ResponseWriter, req *http.Request, sc SharedContext) {
	w.WriteHeader(http.StatusNotFound)
	sc.errorTemplate.Execute(w, resources.ErrorTemplateData{
		Title:       "Not Found",
		Description: "The link you've accessed does not exist",
	})
}

func handler(w http.ResponseWriter, r *http.Request, sc SharedContext) {
	path := r.URL.Path

	err := sc.s.RefetchRoutes()
	if err != nil {
		data := resources.ErrorTemplateData{
			Title:       "Welp, that's not good",
			Description: "There's been an error on our end, please check back later",
		}
		if sc.s.Config.Debug {
			data.Error = err.Error()

			logging.Debug(sc.logger, "error refetching routes: %s", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		sc.errorTemplate.Execute(w, data)
		return
	}

	reply, err := sc.s.LookupRouteUrlFromPath(path)
	if err != nil {
		notFoundHandler(w, r, sc)
		if sc.s.Config.Debug {
			logging.Debug(sc.logger, "Route not found: %s", err)
		}
		return
	}

	http.Redirect(w, r, reply, http.StatusPermanentRedirect)
}
