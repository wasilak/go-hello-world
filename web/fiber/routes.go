package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c *fiber.Ctx) error {
	// Convert fasthttp.Request to http.Request
	req := new(http.Request)
	fasthttpadaptor.ConvertRequest(c.Context(), req, false)
	return c.JSON(s.SetMainResponse(c.Context(), req))
}

func (s *Server) healthRoute(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.JSON(common.HealthResponse{Status: "ok"})
}

func (s *Server) loggerRoute(c *fiber.Ctx) error {
	return c.JSON(s.SetLogLevelResponse(c.Context(), c.Query("level")))
}

func (s *Server) switchRoute(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json")
	return c.JSON(s.SetFrameworkResponse(c.Context(), c.Query("name")))
}
