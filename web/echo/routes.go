package echo

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/wasilak/go-hello-world/web"
)

func mainRoute(c echo.Context) error {
	_, spanResponse := tracer.Start(c.Request().Context(), "response")
	hostname, _ := os.Hostname()
	response := web.APIResponse{
		Host: hostname,
		Request: web.APIResponseRequest{
			Host:       c.Request().Host,
			URL:        c.Request().URL,
			RemoteAddr: c.Request().RemoteAddr,
			RequestURI: c.Request().RequestURI,
			Method:     c.Request().Method,
			Proto:      c.Request().Proto,
			UserAgent:  c.Request().UserAgent(),
			Headers:    c.Request().Header,
		},
	}

	spanResponse.End()
	return c.JSON(http.StatusOK, response)
}

func healthRoute(c echo.Context) error {
	response := web.HealthResponse{Status: "ok"}
	return c.JSON(http.StatusOK, response)
}
