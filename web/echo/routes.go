package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, s.SetMainResponse(c.Request().Context(), c.Request()))
}

func (s *Server) healthRoute(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, common.HealthResponse{Status: "ok"})
}

func (s *Server) loggerRoute(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, s.SetLogLevelResponse(c.Request().Context(), c.QueryParam("level")))
}

func (s *Server) switchRoute(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, s.SetFrameworkResponse(c.Request().Context(), c.QueryParam("name")))
}
