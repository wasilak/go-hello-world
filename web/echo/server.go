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
	slogecho "github.com/samber/slog-echo"
	"github.com/wasilak/go-hello-world/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var logLevel *slog.LevelVar

func Init(ctx context.Context, logLevelConfig *slog.LevelVar, listenAddr *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	tracer = tr
	logLevel = logLevelConfig

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Debug = strings.EqualFold(logLevel.Level().String(), "debug")

	e.Use(slogecho.New(slog.Default()))

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

	e.Use(middleware.Recover())

	e.GET("/", mainRoute)
	e.GET("/health", healthRoute)
	e.GET("/logger", loggerRoute)

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

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	e.Logger.Fatal(e.Start(*listenAddr))
}
