package echo

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/web/common"
	loggergoLib "github.com/wasilak/loggergo/lib"
)

func mainRoute(c echo.Context) error {
	_, spanResponse := tracer.Start(c.Request().Context(), "response")
	response := common.ConstructResponse(c.Request())
	spanResponse.End()
	return c.JSON(http.StatusOK, response)
}

func healthRoute(c echo.Context) error {
	response := common.HealthResponse{Status: "ok"}
	return c.JSON(http.StatusOK, response)
}

func loggerRoute(c echo.Context) error {
	newLogLevelParam := c.QueryParam("level")

	response := common.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergoLib.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(c.Request().Context(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)
	return c.JSON(http.StatusOK, response)
}
