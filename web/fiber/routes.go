package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/wasilak/go-hello-world/web"
)

func mainRoute(c *fiber.Ctx) error {
	_, spanResponse := tracer.Start(c.UserContext(), "response")

	// Convert fasthttp.Request to http.Request
	req := new(http.Request)
	fasthttpadaptor.ConvertRequest(c.Context(), req, false)

	response := web.ConstructResponse(req)
	spanResponse.End()

	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}

func healthRoute(c *fiber.Ctx) error {
	response := web.HealthResponse{Status: "ok"}
	c.Set("Content-Type", "application/json")
	return c.JSON(response)
}
