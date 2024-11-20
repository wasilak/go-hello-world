package echo

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/web"
	"github.com/wasilak/loggergo"
)

func mainRoute(c echo.Context) error {
	_, spanResponse := tracer.Start(c.Request().Context(), "response")
	response := web.ConstructResponse(c.Request())
	spanResponse.End()
	return c.JSON(http.StatusOK, response)
}

func healthRoute(c echo.Context) error {
	response := web.HealthResponse{Status: "ok"}
	return c.JSON(http.StatusOK, response)
}

func loggerRoute(c echo.Context) error {
	newLogLevelParam := c.QueryParam("level")

	response := web.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergo.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(c.Request().Context(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)
	return c.JSON(http.StatusOK, response)
}
