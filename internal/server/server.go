package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/smncd/squishy/internal/filesystem"
)

func New(s *filesystem.SquishyFile) *gin.Engine {
	router := gin.Default()

	router.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")

		err := s.RefreshFile()
		if err != nil {
			c.String(500, "error loading squishyfile")
		}

		reply, ok := s.LookupRoutePath(path)
		if !ok {
			c.String(404, "route not found")
		}

		c.Redirect(http.StatusPermanentRedirect, reply)
	})

	return router
}
