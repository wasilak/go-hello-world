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
)

type Server struct {
	Server *http.Server
	*common.WebServer
}

var p *ginprometheus.Prometheus

type slogWriter struct{}

func (sw slogWriter) Write(p []byte) (n int, err error) {
	slog.Default().Info(string(p))
	return len(p), nil
}

func (s *Server) setup() {
	gin.DefaultWriter = slogWriter{}

	// Create a Gin router
	r := gin.Default()

	// Prometheus Middleware
	if p == nil {
		p = ginprometheus.NewPrometheus(strings.ReplaceAll(utils.GetAppName(), "-", "_"))
	}
	p.Use(r)

	// OpenTelemetry Middleware
	if s.FrameworkOptions.OtelEnabled {
		r.Use(otelgin.Middleware(utils.GetAppName(), otelgin.WithTracerProvider(s.FrameworkOptions.TraceProvider)))
	}

	// Gzip Middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Custom Logging Middleware
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	// Debug Mode
	if strings.EqualFold(s.FrameworkOptions.LogLevelConfig.Level().String(), "debug") {
		gin.SetMode(gin.DebugMode)
	}

	// Define Routes
	r.GET("/", s.mainRoute)
	r.GET("/health", s.healthRoute)
	r.GET("/logger", s.loggerRoute)
	r.GET("/framework", s.switchRoute)
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Optional Statviz
	if s.FrameworkOptions.StatsvizEnabled {
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

	s.Server = &http.Server{
		Addr:    s.FrameworkOptions.ListenAddr,
		Handler: r,
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
		if s.Server == nil {
			s.setup()
		}
		slog.DebugContext(ctx, "Starting server", "address", s.FrameworkOptions.ListenAddr)
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := s.Server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Error stopping web server", "error", err)
	} else {
		slog.InfoContext(ctx, "Web server stopped successfully")
	}

	s.Running = false
}
