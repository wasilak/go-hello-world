package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	otelgotracer "github.com/wasilak/otelgo/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"

	"github.com/wasilak/loggergo"
)

var (
	listenAddr  string
	logLevel    string
	logFormat   string
	otelEnabled bool
)

var tracer = otel.Tracer("go-hello-world")

func main() {

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&logLevel, "log-level", os.Getenv("LOG_LEVEL"), "info")
	flag.StringVar(&logFormat, "log-format", os.Getenv("LOG_FORMAT"), "text")
	flag.BoolVar(&otelEnabled, "otel-enabled", false, "OpenTelemetry traces enabled")
	flag.Parse()

	ctx := context.Background()

	if otelEnabled {
		otelgotracer.InitTracer(ctx, true)
	}

	loggergo.LoggerInit(logLevel, logFormat)

	router := mux.NewRouter()

	router.HandleFunc("/", Chain(rootHandler, Logging()))
	router.HandleFunc("/health", Chain(healthHandler, Logging()))

	router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(statsviz.Ws)
	router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").HandlerFunc(statsviz.Index)

	if otelEnabled {
		router.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	}

	http.ListenAndServe(listenAddr, router)
}
