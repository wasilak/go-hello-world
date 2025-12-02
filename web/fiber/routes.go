package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c *fiber.Ctx) error {
	// Convert fasthttp.Request to http.Request
	req := new(http.Request)
	fasthttpadaptor.ConvertRequest(c.Context(), req, false)
	response := s.SetMainResponse(c.Context(), req)

	if err := c.JSON(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode main response in fiber")
		appErr.AddContext("path", c.Path())
		appErr.LogError(c.Context())
		// Return the original error to be handled by fiber
		return err
	}
	return nil
}

func (s *Server) healthRoute(c *fiber.Ctx) error {
	response := common.HealthResponse{Status: "ok"}

	if err := c.JSON(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode health response in fiber")
		appErr.AddContext("path", c.Path())
		appErr.LogError(c.Context())
		// Return the original error to be handled by fiber
		return err
	}
	return nil
}

func (s *Server) loggerRoute(c *fiber.Ctx) error {
	levelParam := c.Query("level")
	response := s.SetLogLevelResponse(c.Context(), levelParam)

	if err := c.JSON(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode logger response in fiber")
		appErr.AddContext("path", c.Path())
		appErr.AddContext("log_level", levelParam)
		appErr.LogError(c.Context())
		// Return the original error to be handled by fiber
		return err
	}
	return nil
}

func (s *Server) switchRoute(c *fiber.Ctx) error {
	nameParam := c.Query("name")
	response := s.SetFrameworkResponse(c.Context(), nameParam)

	c.Set("Content-Type", "application/json")
	if err := c.JSON(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode switch response in fiber")
		appErr.AddContext("path", c.Path())
		appErr.AddContext("framework", nameParam)
		appErr.LogError(c.Context())
		// Return the original error to be handled by fiber
		return err
	}
	return nil
}
