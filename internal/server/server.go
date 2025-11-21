package server

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/config"
	"gitlab.com/smncd/squishy/internal/logging"
	"gitlab.com/smncd/squishy/internal/router"
)

type SharedContext struct {
	cfg    *config.Config
	routes *config.Routes
	logger *log.Logger
}

func New(cfg *config.Config, routes *config.Routes, logger *log.Logger) *http.Server {
	sc := SharedContext{
		cfg:    cfg,
		routes: routes,
		logger: logger,
	}

	router := router.New(sc, logger)

	router.NoRoute(notFoundHandler)

	router.GET("/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		Handler: router,
	}

	return server
}

func notFoundHandler(w http.ResponseWriter, req *http.Request, sc SharedContext) {
	w.WriteHeader(http.StatusNotFound)
}

func handler(w http.ResponseWriter, r *http.Request, sc SharedContext) {
	path := r.URL.Path

	err := sc.routes.Refetch()
	if err != nil {
		if sc.cfg.Debug {
			logging.Debug(sc.logger, "error refetching routes: %s", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reply, err := sc.routes.LookupUrlFromPath(path)
	if err != nil {
		notFoundHandler(w, r, sc)
		if sc.cfg.Debug {
			logging.Debug(sc.logger, "Route not found: %s", err)
		}
		return
	}

	http.Redirect(w, r, reply, http.StatusPermanentRedirect)
}
