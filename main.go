package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"log/slog"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	otelgotracer "github.com/wasilak/otelgo/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"

	"github.com/wasilak/loggergo"
)

var (
	listenAddr             string
	logLevel               string
	logFormat              string
	otelEnabled            bool
	otelHostMetricsEnabled bool
)

var tracer = otel.Tracer("go-hello-world")

func main() {

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&logLevel, "log-level", os.Getenv("LOG_LEVEL"), "info")
	flag.StringVar(&logFormat, "log-format", os.Getenv("LOG_FORMAT"), "text")
	flag.BoolVar(&otelEnabled, "otel-enabled", false, "OpenTelemetry traces enabled")
	flag.BoolVar(&otelHostMetricsEnabled, "otel-host-metrics", false, "OpenTelemetry host metrics enabled")
	flag.Parse()

	ctx := context.Background()

	if otelEnabled {
		otelGoTracingConfig := otelgotracer.OtelGoTracingConfig{
			HostMetricsEnabled: false,
		}
		err := otelgotracer.InitTracer(ctx, otelGoTracingConfig)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			os.Exit(1)
		}
	}

	loggerConfig := loggergo.LoggerGoConfig{
		Level:  logLevel,
		Format: logFormat,
	}

	_, err := loggergo.LoggerInit(loggerConfig)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", Chain(rootHandler, Logging()))
	router.HandleFunc("/health", Chain(healthHandler, Logging()))

	// Create statsviz server and register the handlers on the router.
	srv, _ := statsviz.NewServer()
	router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(srv.Ws())
	router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(srv.Index())

	if otelEnabled {
		router.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	}

	http.ListenAndServe(listenAddr, router)
}
