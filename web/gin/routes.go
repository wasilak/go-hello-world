package gin

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wasilak/go-hello-world/web/common"
	loggergoLib "github.com/wasilak/loggergo/lib"
)

func mainRoute(c *gin.Context) {
	_, spanResponse := tracer.Start(c.Request.Context(), "response")
	response := common.ConstructResponse(c.Request)
	spanResponse.End()
	c.JSON(http.StatusOK, response)
}

func healthRoute(c *gin.Context) {
	response := common.HealthResponse{Status: "ok"}
	c.JSON(http.StatusOK, response)
}

func loggerRoute(c *gin.Context) {
	newLogLevelParam := c.Query("level")

	response := common.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergoLib.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(c.Request.Context(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)
	c.JSON(http.StatusOK, response)
}
