package gorilla

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sloghttp "github.com/samber/slog-http"
	"github.com/wasilak/go-hello-world/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var logLevel *slog.LevelVar

func Init(ctx context.Context, logLevelConfig *slog.LevelVar, listenAddr *string, otelEnabled, statsvizEnabled *bool, tr trace.Tracer) {
	tracer = tr
	logLevel = logLevelConfig
	router := mux.NewRouter()

	// Prometheus middleware and metrics endpoint
	router.Use(prometheusMiddleware)
	router.Path("/metrics").Handler(promhttp.Handler())

	// Application-specific routes
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/logger", loggerHandler)

	if *statsvizEnabled {
		// Create statsviz server and register the handlers on the router
		srv, _ := statsviz.NewServer()

		slog.DebugContext(ctx, "Statsviz enabled", "address", "/debug/statsviz/")
		router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(srv.Ws())
		router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(srv.Index())
	}

	if *otelEnabled {
		router.Use(otelmux.Middleware(utils.GetAppName()))
	}

	// Wrap the router with sloghttp middleware
	handler := sloghttp.Recovery(router)            // Recovery middleware
	handler = sloghttp.New(slog.Default())(handler) // Logging middleware

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	http.ListenAndServe(*listenAddr, handler) // Use the wrapped handler
}
