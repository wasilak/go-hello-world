package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c echo.Context) error {
	ctx := c.Request().Context()
	response := s.SetMainResponse(ctx, c.Request())

	if err := c.JSON(http.StatusOK, response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode main response in echo")
		appErr.AddContext("path", c.Path())
		appErr.LogError(ctx)
		// Return the original error to be handled by echo
		return err
	}
	return nil
}

func (s *Server) healthRoute(c echo.Context) error {
	ctx := c.Request().Context()
	response := common.HealthResponse{Status: "ok"}

	if err := c.JSON(http.StatusOK, response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode health response in echo")
		appErr.AddContext("path", c.Path())
		appErr.LogError(ctx)
		// Return the original error to be handled by echo
		return err
	}
	return nil
}

func (s *Server) loggerRoute(c echo.Context) error {
	ctx := c.Request().Context()
	levelParam := c.QueryParam("level")
	response := s.SetLogLevelResponse(ctx, levelParam)

	if err := c.JSON(http.StatusOK, response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode logger response in echo")
		appErr.AddContext("path", c.Path())
		appErr.AddContext("log_level", levelParam)
		appErr.LogError(ctx)
		// Return the original error to be handled by echo
		return err
	}
	return nil
}

func (s *Server) switchRoute(c echo.Context) error {
	ctx := c.Request().Context()
	nameParam := c.QueryParam("name")
	response := s.SetFrameworkResponse(ctx, nameParam)

	if err := c.JSON(http.StatusOK, response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode switch response in echo")
		appErr.AddContext("path", c.Path())
		appErr.AddContext("framework", nameParam)
		appErr.LogError(ctx)
		// Return the original error to be handled by echo
		return err
	}
	return nil
}
