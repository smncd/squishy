package server

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/templates"
)

//go:embed static/*
var staticFS embed.FS

func New(s *filesystem.SquishyFile) *http.Server {
	if !s.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.StaticFS("/static", http.FS(staticFS))

	f := template.Must(template.ParseFS(templates.FS, "*.html"))
	router.SetHTMLTemplate(f)

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		err := s.RefetchRoutes()
		if err != nil {
			data := templates.ErrorPageData{
				Title:       "Welp, that's not good",
				Description: "There's been an error on our end, please check back later",
			}
			if s.Config.Debug {
				data.Error = err.Error()
			}

			c.HTML(http.StatusInternalServerError, "error.html", data)
			return
		}

		reply, ok := s.LookupRoutePath(path)
		if !ok {
			c.HTML(http.StatusNotFound, "error.html", templates.ErrorPageData{
				Title:       "Not Found",
				Description: "The link you've accessed does not exist",
			})
			return
		}

		c.Redirect(http.StatusPermanentRedirect, reply)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router.Handler(),
	}

	return srv
}
