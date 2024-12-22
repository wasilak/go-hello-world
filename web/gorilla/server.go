package gorilla

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sloghttp "github.com/samber/slog-http"
	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type Server struct {
	Server *http.Server
	wg     sync.WaitGroup
	*common.WebServer
}

func (s *Server) setup(ctx context.Context) {
	router := mux.NewRouter()

	// Prometheus middleware and metrics endpoint
	router.Use(prometheusMiddleware)
	router.Path("/metrics").Handler(promhttp.Handler())

	// Application-specific routes
	router.HandleFunc("/", s.rootHandler)
	router.HandleFunc("/health", s.healthHandler)
	router.HandleFunc("/logger", s.loggerHandler)
	router.HandleFunc("/framework", s.switchRoute)

	if s.FrameworkOptions.StatsvizEnabled {
		// Create statsviz server and register the handlers on the router
		srv, _ := statsviz.NewServer()

		slog.DebugContext(ctx, "Statsviz enabled", "address", "/debug/statsviz/")
		router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(srv.Ws())
		router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(srv.Index())
	}

	if s.FrameworkOptions.OtelEnabled {
		router.Use(otelmux.Middleware(utils.GetAppName()))
	}

	// Wrap the router with sloghttp middleware
	handler := sloghttp.Recovery(router)            // Recovery middleware
	handler = sloghttp.New(slog.Default())(handler) // Logging middleware

	s.Server = &http.Server{
		Addr:    s.FrameworkOptions.ListenAddr,
		Handler: handler,
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
			s.setup(ctx)
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
