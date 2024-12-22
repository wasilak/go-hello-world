package fiber

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/arl/statsviz"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"

	slogfiber "github.com/samber/slog-fiber"
)

type Server struct {
	Server *fiber.App
	*common.WebServer
}

func (s *Server) setup(ctx context.Context) {

	// Initialize Fiber app
	s.Server = fiber.New(fiber.Config{
		DisableStartupMessage: true, // Disable the Fiber banner
	})

	// Prometheus Middleware
	prometheus := fiberprometheus.New(utils.GetAppName())
	prometheus.RegisterAt(s.Server, "/metrics")
	s.Server.Use(prometheus.Middleware)

	// OpenTelemetry Middleware
	if s.FrameworkOptions.OtelEnabled {
		s.Server.Use(otelfiber.Middleware())
	}

	// Gzip Middleware
	s.Server.Use(compress.New())

	// Custom Logging Middleware
	s.Server.Use(slogfiber.New(slog.Default()))

	// Define Routes
	s.Server.Get("/", s.mainRoute)
	s.Server.Get("/health", s.healthRoute)
	s.Server.Get("/logger", s.loggerRoute)
	s.Server.Get("/framework", s.switchRoute)

	// Optional Statviz
	if s.FrameworkOptions.StatsvizEnabled {
		mux := http.NewServeMux()

		// Register statsviz handlerson the mux.
		statsviz.Register(mux)

		// Register Statsviz routes on the Fiber app
		s.Server.Use("/debug/statsviz", adaptor.HTTPHandler(mux))
		s.Server.Get("/debug/statsviz/*", adaptor.HTTPHandler(mux))
	}

	slog.DebugContext(ctx, "Starting server", "address", s.FrameworkOptions.ListenAddr)
}

func (s *Server) Start(ctx context.Context) {
	s.MU.Lock()
	defer s.MU.Unlock()

	if s.Running {
		slog.DebugContext(ctx, "Web server is already running")
		return
	}

	go func() {
		s.setup(ctx)
		slog.DebugContext(ctx, "Starting server", "address", s.FrameworkOptions.ListenAddr)
		if err := s.Server.Listen(s.FrameworkOptions.ListenAddr); err != nil {
			slog.ErrorContext(ctx, "Server exited with error", "error", err)
			os.Exit(1)
		}
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
	// shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	if err := s.Server.Shutdown(); err != nil {
		slog.ErrorContext(ctx, "Error stopping web server", "error", err)
	} else {
		slog.InfoContext(ctx, "Web server stopped successfully")
	}

	s.Running = false
}
