package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/smncd/squishy/internal/embedfs"
	"gitlab.com/smncd/squishy/internal/filesystem"
)

func New(s *filesystem.SquishyFile) *http.Server {
	if !s.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	f := template.Must(template.ParseFS(embedfs.FS, "*.html"))
	router.SetHTMLTemplate(f)

	router.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")

		err := s.RefetchRoutes()
		if err != nil {
			errorMessage := "error fetching routes"
			if s.Config.Debug {
				errorMessage = err.Error()
			}
			c.String(http.StatusInternalServerError, errorMessage)
		}

		reply, ok := s.LookupRoutePath(path)
		if !ok {
			c.HTML(http.StatusNotFound, "404.html", nil)
		}

		c.Redirect(http.StatusPermanentRedirect, reply)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router.Handler(),
	}

	return srv
}
