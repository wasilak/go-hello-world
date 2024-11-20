package gin

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wasilak/go-hello-world/web"
	"github.com/wasilak/loggergo"
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

func loggerRoute(c *gin.Context) {
	newLogLevelParam := c.Query("level")

	response := web.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergo.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(c.Request.Context(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)
	c.JSON(http.StatusOK, response)
}
