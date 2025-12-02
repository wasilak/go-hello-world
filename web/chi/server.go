package chi

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	slogchi "github.com/samber/slog-chi"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
)

type Server struct {
	Server *http.Server
	wg     sync.WaitGroup
	*common.WebServer
}

var promMiddleware func(next http.Handler) http.Handler

func (s *Server) setup() {

	r := chi.NewRouter()

	// Register Go runtime metrics only if not already registered
	goCollector := collectors.NewGoCollector()
	processCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})

	// Use shared utility to prevent duplicate registration
	common.RegisterCollectorIfNotRegistered(goCollector)
	common.RegisterCollectorIfNotRegistered(processCollector)

	if promMiddleware == nil {
		promMiddleware = chiprometheus.NewMiddleware(strings.ReplaceAll(utils.GetAppName(), "-", "_"))
	}
	r.Use(promMiddleware)

	// OpenTelemetry Middleware
	if s.FrameworkOptions.OtelEnabled {
		r.Use(otelchi.Middleware(utils.GetAppName(), otelchi.WithFilter(func(r *http.Request) bool {
			return !strings.Contains(r.URL.Path, "public/dist") && !strings.Contains(r.URL.Path, "health")
		})))
	}

	// Gzip Middleware
	r.Use(middleware.NewCompressor(5).Handler)

	// Custom Logging Middleware
	r.Use(slogchi.New(slog.Default()))

	// Define Routes
	r.Get("/", s.mainRoute)
	r.Get("/health", s.healthRoute)
	r.Get("/logger", s.loggerRoute)
	r.Get("/framework", s.switchRoute)
	r.Handle("/metrics", promhttp.Handler())

	// Optional Statviz
	if s.FrameworkOptions.StatsvizEnabled {
		// Create statsviz server.
		srv, _ := statsviz.NewServer()

		r.Get("/debug/statsviz/ws", srv.Ws())
		r.Get("/debug/statsviz", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/debug/statsviz/", http.StatusMovedPermanently)
		})
		r.Handle("/debug/statsviz/*", srv.Index())
	}

	s.Server = &http.Server{
		Addr:    s.FrameworkOptions.ListenAddr,
		Handler: r,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.MU.Lock()
	defer s.MU.Unlock()

	s.wg.Add(1)

	if s.Running {
		slog.DebugContext(ctx, "Web server is already running")
		return
	}

	go func() {
		defer s.wg.Done()

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
