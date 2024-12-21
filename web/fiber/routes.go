package fiber

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/wasilak/go-hello-world/web/common"
	loggergoLib "github.com/wasilak/loggergo/lib"
)

func mainRoute(c *fiber.Ctx) error {
	_, spanResponse := tracer.Start(c.UserContext(), "response")

	// Convert fasthttp.Request to http.Request
	req := new(http.Request)
	fasthttpadaptor.ConvertRequest(c.Context(), req, false)

	response := common.ConstructResponse(req)
	spanResponse.End()

	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}

func healthRoute(c *fiber.Ctx) error {
	response := common.HealthResponse{Status: "ok"}
	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}

func loggerRoute(c *fiber.Ctx) error {
	newLogLevelParam := c.Query("level")

	response := common.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergoLib.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(c.UserContext(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)

	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}
