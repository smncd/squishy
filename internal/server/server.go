package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/smncd/squishy/internal/filesystem"
)

func New(s *filesystem.SquishyFile) *http.Server {
	if !s.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")

		err := s.RefetchFile()
		if err != nil {
			c.String(500, "error loading squishyfile")
		}

		reply, ok := s.LookupRoutePath(path)
		if !ok {
			c.String(404, "route not found")
		}

		c.Redirect(http.StatusPermanentRedirect, reply)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port),
		Handler: router.Handler(),
	}

	return srv
}
