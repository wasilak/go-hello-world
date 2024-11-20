package gin

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/arl/statsviz"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/wasilak/go-hello-world/utils"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var logLevel *slog.LevelVar

type slogWriter struct{}

func (sw slogWriter) Write(p []byte) (n int, err error) {
	slog.Default().Info(string(p))
	return len(p), nil
}

func Init(ctx context.Context, logLevelConfig *slog.LevelVar, listenAddr *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer, traceProvider *trace.TracerProvider) {
	tracer = tr
	logLevel = logLevelConfig

	gin.DefaultWriter = slogWriter{}

	// Create a Gin router
	r := gin.Default()

	// Prometheus Middleware
	p := ginprometheus.NewPrometheus(strings.ReplaceAll(utils.GetAppName(), "-", "_"))
	p.Use(r)

	// OpenTelemetry Middleware
	if *otelEnabled {
		r.Use(otelgin.Middleware(utils.GetAppName(), otelgin.WithTracerProvider(*traceProvider)))
	}

	// Gzip Middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Custom Logging Middleware
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	// Debug Mode
	if strings.EqualFold(logLevel.Level().String(), "debug") {
		gin.SetMode(gin.DebugMode)
		slog.DebugContext(ctx, "Debug mode enabled")
	}

	// Define Routes
	r.GET("/", mainRoute)
	r.GET("/health", healthRoute)
	r.GET("/logger", loggerRoute)
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Optional Statviz
	if *statsvizEnabled {
		// Create statsviz server.
		srv, _ := statsviz.NewServer()

		// Specific route for WebSocket
		r.GET("/debug/statsviz/ws", gin.WrapF(srv.Ws()))

		// Static route for the index
		r.GET("/debug/statsviz", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/debug/statsviz/")
		})

		// Define more specific routes for statsviz
		r.GET("/debug/statsviz/plots", gin.WrapH(srv.Index()))
		r.GET("/debug/statsviz/metrics", gin.WrapH(srv.Index()))
		// Add other specific paths as needed
	}

	// Start Server
	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	if err := r.Run(*listenAddr); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", "error", err)
		os.Exit(1)
	}
}
