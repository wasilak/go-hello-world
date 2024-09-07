package echo

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/arl/statsviz"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	slogecho "github.com/samber/slog-echo"
	"github.com/wasilak/go-hello-world/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Init(ctx context.Context, listenAddr, logLevel *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	tracer = tr

	e := echo.New()

	if *otelEnabled {
		e.Use(otelecho.Middleware(utils.GetAppName(), otelecho.WithSkipper(func(c echo.Context) bool {
			return strings.Contains(c.Path(), "public/dist") || strings.Contains(c.Path(), "health")
		})))
	}

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "metrics")
		},
	}))

	e.Use(echoprometheus.NewMiddleware(strings.ReplaceAll(utils.GetAppName(), "-", "_")))

	e.HideBanner = true
	e.HidePort = true

	if strings.EqualFold(*logLevel, "debug") {
		e.Logger.SetLevel(log.DEBUG)
		e.Debug = true
	}

	e.Use(slogecho.New(slog.Default()))
	e.Use(middleware.Recover())

	e.GET("/", mainRoute)
	e.GET("/health", healthRoute)

	e.GET("/metrics", echoprometheus.NewHandler())

	if *statsvizEnabled {
		// Create statsviz server and register the handlers on the router.
		mux := http.NewServeMux()

		// Register statsviz handlerson the mux.
		statsviz.Register(mux)

		// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
		e.GET("/debug/statsviz/", echo.WrapHandler(mux))
		// Serve static content for statsviz UI
		e.GET("/debug/statsviz/*", echo.WrapHandler(mux))
	}

	slog.DebugContext(ctx, "Features supported", "loggergo", true, "statsviz", true, "tracing", true)

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	e.Logger.Fatal(e.Start(*listenAddr))
}
