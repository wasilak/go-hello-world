package gorilla

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Init(ctx context.Context, listenAddr *string, otelEnabled, profilingEnabled *bool, t trace.Tracer) {
	tracer = t
	router := mux.NewRouter()

	router.HandleFunc("/", chain(rootHandler, logging()))
	router.HandleFunc("/health", chain(healthHandler, logging()))

	if *profilingEnabled {
		// Create statsviz server and register the handlers on the router.
		srv, _ := statsviz.NewServer()
		router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(srv.Ws())
		router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(srv.Index())
	}

	if *otelEnabled {
		router.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	}

	slog.DebugContext(ctx, "Features supported", "loggergo", true, "profiling", true, "tracing", true)

	slog.DebugContext(ctx, "Starting server", "address", *listenAddr)
	http.ListenAndServe(*listenAddr, router)
}
