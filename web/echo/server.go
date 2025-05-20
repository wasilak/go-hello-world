package echo

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/arl/statsviz"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	slogecho "github.com/samber/slog-echo"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	Server *echo.Echo
	*common.WebServer
}

func (s *Server) setup() {
	s.Server = echo.New()

	s.Server.HideBanner = true
	s.Server.HidePort = true

	s.Server.Debug = strings.EqualFold(s.FrameworkOptions.LogLevelConfig.Level().String(), "debug")

	if s.FrameworkOptions.OtelEnabled {
		s.Server.Use(otelecho.Middleware(utils.GetAppName(), otelecho.WithTracerProvider(s.FrameworkOptions.TraceProvider), otelecho.WithSkipper(func(c echo.Context) bool {
			return strings.Contains(c.Path(), "public/dist") || strings.Contains(c.Path(), "health")
		})))
	}

	s.Server.Use(slogecho.New(slog.Default()))

	s.Server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "metrics")
		},
	}))

	echoprometheusConfig := echoprometheus.MiddlewareConfig{
		Subsystem:  strings.ReplaceAll(utils.GetAppName(), "-", "_"),
		Registerer: prometheus.Registerer(prometheus.NewRegistry()),
	}
	s.Server.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheusConfig))

	s.Server.Use(middleware.Recover())

	s.Server.GET("/", s.mainRoute)
	s.Server.GET("/health", s.healthRoute)
	s.Server.GET("/logger", s.loggerRoute)
	s.Server.GET("/framework", s.switchRoute)

	s.Server.GET("/metrics", echoprometheus.NewHandler())

	if s.FrameworkOptions.StatsvizEnabled {
		// Create statsviz server and register the handlers on the router.
		mux := http.NewServeMux()

		// Register statsviz handlerson the mux.
		statsviz.Register(mux)

		// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
		s.Server.GET("/debug/statsviz/", echo.WrapHandler(mux))
		// Serve static content for statsviz UI
		s.Server.GET("/debug/statsviz/*", echo.WrapHandler(mux))
	}
}

func (s *Server) Start(ctx context.Context) {
	s.MU.Lock()
	defer s.MU.Unlock()

	if s.Running {
		slog.DebugContext(ctx, "Web server is already running")
		return
	}

	go func() {
		s.setup()
		slog.DebugContext(ctx, "Starting server", "address", s.FrameworkOptions.ListenAddr)
		s.Server.Start(s.FrameworkOptions.ListenAddr)
	}()

	s.Running = true
}

// Stop gracefully stops the web server.
func (s *Server) Stop(ctx context.Context) {
	s.MU.Lock()
	defer s.MU.Unlock()

	if !s.Running {
		slog.DebugContext(ctx, "Web server is not running")
		return
	}

	slog.InfoContext(ctx, "Stopping web server")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Error stopping web server", "error", err)
	} else {
		slog.InfoContext(ctx, "Web server stopped successfully")
	}

	s.Running = false
}
