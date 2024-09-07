package chi

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/arl/statsviz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	slogchi "github.com/samber/slog-chi"
	"github.com/wasilak/go-hello-world/utils"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Init(ctx context.Context, listenAddr, logLevel *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	slog.DebugContext(ctx, "Features supported", "loggergo", true, "statsviz", true, "tracing", true)
	tracer = tr

	r := chi.NewRouter()

	r.Use(chiprometheus.NewMiddleware(strings.ReplaceAll(utils.GetAppName(), "-", "_")))

	// OpenTelemetry Middleware
	if *otelEnabled {
		r.Use(otelchi.Middleware(utils.GetAppName(), otelchi.WithFilter(func(r *http.Request) bool {
			return !strings.Contains(r.URL.Path, "public/dist") && !strings.Contains(r.URL.Path, "health")
		})))
	}

	// Gzip Middleware
	r.Use(middleware.NewCompressor(5).Handler)

	// Custom Logging Middleware
	r.Use(slogchi.New(slog.Default()))

	// Debug Mode
	if strings.EqualFold(*logLevel, "debug") {
		slog.DebugContext(ctx, "Debug mode enabled")
	}

	// Define Routes
	r.Get("/", mainRoute)
	r.Get("/health", healthRoute)
	r.Handle("/metrics", promhttp.Handler())

	// Optional Statviz
	if *statsvizEnabled {
		// Create statsviz server.
		srv, _ := statsviz.NewServer()

		r.Get("/debug/statsviz/ws", srv.Ws())
		r.Get("/debug/statsviz", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/debug/statsviz/", http.StatusMovedPermanently)
		})
		r.Handle("/debug/statsviz/*", srv.Index())
	}

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, r); err != nil {
		slog.ErrorContext(ctx, "Server exited with error", "error", err)
		os.Exit(1)
	}
}
