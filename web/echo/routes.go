package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/web"
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
