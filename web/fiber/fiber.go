package fiber

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/arl/statsviz"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/wasilak/go-hello-world/utils"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Init(ctx context.Context, listenAddr, logLevel *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	slog.DebugContext(ctx, "Features supported", "loggergo", true, "statsviz", true, "tracing", true)
	tracer = tr

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true, // Disable the Fiber banner
	})

	// Prometheus Middleware
	prometheus := fiberprometheus.New(utils.GetAppName())
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// OpenTelemetry Middleware
	if *otelEnabled {
		app.Use(otelfiber.Middleware())
	}

	// Gzip Middleware
	app.Use(compress.New())

	// Custom Logging Middleware
	app.Use(func(c *fiber.Ctx) error {
		slog.InfoContext(ctx, "Incoming request", "method", c.Method(), "path", c.Path())
		return c.Next()
	})

	// Debug Mode
	if strings.EqualFold(*logLevel, "debug") {
		slog.DebugContext(ctx, "Debug mode enabled")
	}

	// Define Routes
	app.Get("/", func(c *fiber.Ctx) error { return mainRoute(c) })
	app.Get("/health", func(c *fiber.Ctx) error { return healthRoute(c) })

	// Optional Statviz
	if *statsvizEnabled {
		ws := http.NewServeMux()

		// Create statsviz server.
		srv, err := statsviz.NewServer()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create statsviz server", "error", err)
			os.Exit(1)
		}

		// Register Statsviz server on the fasthttp router.
		app.Use("/debug/statsviz", adaptor.HTTPHandler(srv.Index()))
		ws.HandleFunc("/debug/statsviz/ws", srv.Ws())
	}

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	if err := app.Listen(*listenAddr); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", "error", err)
		os.Exit(1)
	}
}
