package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c *gin.Context) {
	ctx := c.Request.Context()
	response := s.SetMainResponse(ctx, c.Request)

	// Gin automatically handles JSON marshaling errors internally
	c.JSON(http.StatusOK, response)
}

func (s *Server) healthRoute(c *gin.Context) {
	// Gin automatically handles JSON marshaling errors internally
	c.JSON(http.StatusOK, common.HealthResponse{Status: "ok"})
}

func (s *Server) loggerRoute(c *gin.Context) {
	ctx := c.Request.Context()
	levelParam := c.Query("level")
	response := s.SetLogLevelResponse(ctx, levelParam)

	// Gin automatically handles JSON marshaling errors internally
	c.JSON(http.StatusOK, response)
}

func (s *Server) switchRoute(c *gin.Context) {
	ctx := c.Request.Context()
	nameParam := c.Query("name")
	response := s.SetFrameworkResponse(ctx, nameParam)

	// Gin automatically handles JSON marshaling errors internally
	c.JSON(http.StatusOK, response)
}
