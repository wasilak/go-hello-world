package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c *gin.Context) {
	c.JSON(http.StatusOK, s.SetMainResponse(c.Request.Context(), c.Request))
}

func (s *Server) healthRoute(c *gin.Context) {
	response := common.HealthResponse{Status: "ok"}
	c.JSON(http.StatusOK, response)
}

func (s *Server) loggerRoute(c *gin.Context) {
	c.JSON(http.StatusOK, s.SetLogLevelResponse(c.Request.Context(), c.Query("level")))
}

func (s *Server) switchRoute(c *gin.Context) {
	c.JSON(http.StatusOK, s.SetFrameworkResponse(c.Request.Context(), c.Query("name")))
}
