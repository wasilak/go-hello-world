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
	"github.com/wasilak/go-hello-world/web/common"
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

func Init(ctx context.Context, frameworkOptions common.FrameworkOptions) {
	tracer = frameworkOptions.Tracer
	logLevel = frameworkOptions.LogLevelConfig

	gin.DefaultWriter = slogWriter{}

	// Create a Gin router
	r := gin.Default()

	// Prometheus Middleware
	p := ginprometheus.NewPrometheus(strings.ReplaceAll(utils.GetAppName(), "-", "_"))
	p.Use(r)

	// OpenTelemetry Middleware
	if frameworkOptions.OtelEnabled {
		r.Use(otelgin.Middleware(utils.GetAppName(), otelgin.WithTracerProvider(frameworkOptions.TraceProvider)))
	}

	// Gzip Middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Custom Logging Middleware
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	// Debug Mode
	if strings.EqualFold(logLevel.Level().String(), "debug") {
		gin.SetMode(gin.DebugMode)
	}

	// Define Routes
	r.GET("/", mainRoute)
	r.GET("/health", healthRoute)
	r.GET("/logger", loggerRoute)
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Optional Statviz
	if frameworkOptions.StatsvizEnabled {
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
	slog.DebugContext(ctx, "Starting server", "address", frameworkOptions.ListenAddr)
	if err := r.Run(frameworkOptions.ListenAddr); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", "error", err)
		os.Exit(1)
	}
}
