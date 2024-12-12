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

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"go.opentelemetry.io/otel/trace"

	slogfiber "github.com/samber/slog-fiber"
)

var tracer trace.Tracer
var logLevel *slog.LevelVar

func Init(ctx context.Context, logLevelConfig *slog.LevelVar, listenAddr *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	tracer = tr
	logLevel = logLevelConfig

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
	app.Use(slogfiber.New(slog.Default()))

	// Define Routes
	app.Get("/", func(c *fiber.Ctx) error { return mainRoute(c) })
	app.Get("/health", func(c *fiber.Ctx) error { return healthRoute(c) })
	app.Get("/logger", func(c *fiber.Ctx) error { return loggerRoute(c) })

	// Optional Statviz
	if *statsvizEnabled {
		mux := http.NewServeMux()

		// Register statsviz handlerson the mux.
		statsviz.Register(mux)

		// Register Statsviz routes on the Fiber app
		app.Use("/debug/statsviz", adaptor.HTTPHandler(mux))
		app.Get("/debug/statsviz/*", adaptor.HTTPHandler(mux))
	}

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)

	if err := app.Listen(*listenAddr); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", "error", err)
		os.Exit(1)
	}
}
