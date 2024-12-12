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
var logLevel *slog.LevelVar

func Init(ctx context.Context, logLevelConfig *slog.LevelVar, listenAddr *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	tracer = tr
	logLevel = logLevelConfig

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

	// Define Routes
	r.Get("/", mainRoute)
	r.Get("/health", healthRoute)
	r.Get("/logger", loggerRoute)
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
