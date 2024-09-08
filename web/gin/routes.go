package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wasilak/go-hello-world/web"
)

func mainRoute(c *gin.Context) {
	_, spanResponse := tracer.Start(c.Request.Context(), "response")
	response := web.ConstructResponse(c.Request)
	spanResponse.End()
	c.JSON(http.StatusOK, response)
}

func healthRoute(c *gin.Context) {
	response := web.HealthResponse{Status: "ok"}
	c.JSON(http.StatusOK, response)
}
