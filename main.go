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
	"github.com/wasilak/profilego"
)

var tracer = otel.Tracer(GetAppName())

func main() {

	listenAddr := flag.String("listen-addr", "127.0.0.1:5000", "server listen address")
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "log level (debug, info, warn, error, fatal)")
	logFormat := flag.String("log-format", os.Getenv("LOG_FORMAT"), "log format (json, plain, otel)")
	devMode := flag.Bool("dev-mode", false, "Development mode")
	otelEnabled := flag.Bool("otel-enabled", false, "OpenTelemetry traces enabled")
	otelHostMetricsEnabled := flag.Bool("otel-host-metrics", false, "OpenTelemetry host metrics enabled")
	otelRuntimeMetricsEnabled := flag.Bool("otel-runtime-metrics", false, "OpenTelemetry runtime metrics enabled")
	profilingEnabled := flag.Bool("profiling-enabled", false, "Profiling enabled")
	profilingAddress := flag.String("profiling-address", "127.0.0.1:4040", "Profiling address")
	flag.Parse()

	if *profilingEnabled {
		profileGoConfig := profilego.Config{
			ApplicationName: GetAppName(),
			ServerAddress:   *profilingAddress,
			Type:            "pyroscope",
			Tags:            map[string]string{},
		}
		profilego.Init(profileGoConfig)
	}

	ctx := context.Background()

	loggerConfig := loggergo.Config{
		Level:        loggergo.LogLevelFromString(*logLevel),
		Format:       loggergo.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      *devMode,
		Output:       loggergo.OutputConsole,
	}

	if *otelEnabled {
		otelGoTracingConfig := otelgotracer.Config{
			HostMetricsEnabled:    *otelHostMetricsEnabled,
			RuntimeMetricsEnabled: *otelRuntimeMetricsEnabled,
		}
		_, _, err := otelgotracer.Init(ctx, otelGoTracingConfig)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			os.Exit(1)
		}

		loggerConfig.OtelServiceName = os.Getenv("OTEL_SERVICE_NAME")
		loggerConfig.Output = loggergo.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/go-hello-world"
		loggerConfig.OtelTracingEnabled = false
	}

	_, err := loggergo.LoggerInit(ctx, loggerConfig)
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

	if *otelEnabled {
		router.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	}

	http.ListenAndServe(*listenAddr, router)
}
