package main

import (
	"context"
	"flag"
	"os"

	"log/slog"

	otelgotracer "github.com/wasilak/otelgo/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/chi"
	"github.com/wasilak/go-hello-world/web/echo"
	"github.com/wasilak/go-hello-world/web/gin"
	"github.com/wasilak/go-hello-world/web/gorilla"
	"github.com/wasilak/loggergo"
	"github.com/wasilak/profilego"
)

var tracer = otel.Tracer(utils.GetAppName())

func main() {
	ctx := context.Background()

	listenAddr := flag.String("listen-addr", "127.0.0.1:3000", "server listen address")
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "log level (debug, info, warn, error, fatal)")
	logFormat := flag.String("log-format", os.Getenv("LOG_FORMAT"), "log format (json, plain, otel)")
	otelEnabled := flag.Bool("otel-enabled", false, "OpenTelemetry traces enabled")
	otelHostMetricsEnabled := flag.Bool("otel-host-metrics", false, "OpenTelemetry host metrics enabled")
	otelRuntimeMetricsEnabled := flag.Bool("otel-runtime-metrics", false, "OpenTelemetry runtime metrics enabled")
	statsvizEnabled := flag.Bool("statsviz-enabled", false, "statsviz enabled")
	profilingEnabled := flag.Bool("profiling-enabled", false, "Profiling enabled")
	profilingAddress := flag.String("profiling-address", "127.0.0.1:4040", "Profiling address")
	webFramework := flag.String("web-framework", "gorilla", "Web framework (gorilla, echo, gin, chi)")
	flag.Parse()

	if *profilingEnabled {
		profileGoConfig := profilego.Config{
			ApplicationName: utils.GetAppName(),
			ServerAddress:   *profilingAddress,
			Type:            "pyroscope",
			Tags:            map[string]string{},
		}
		profilego.Init(profileGoConfig)
	}

	loggerConfig := loggergo.Config{
		Level:        loggergo.LogLevelFromString(*logLevel),
		Format:       loggergo.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      loggergo.LogLevelFromString(*logLevel) == slog.LevelDebug && *logFormat == "plain",
		Output:       loggergo.OutputConsole,
	}

	var traceProvider trace.TracerProvider
	var err error

	if *otelEnabled {
		otelGoTracingConfig := otelgotracer.Config{
			HostMetricsEnabled:    *otelHostMetricsEnabled,
			RuntimeMetricsEnabled: *otelRuntimeMetricsEnabled,
		}
		_, traceProvider, err = otelgotracer.Init(ctx, otelGoTracingConfig)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			os.Exit(1)
		}

		loggerConfig.OtelServiceName = utils.GetAppName()
		loggerConfig.Output = loggergo.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/go-hello-world"
		loggerConfig.OtelTracingEnabled = false
	}

	_, err = loggergo.LoggerInit(ctx, loggerConfig)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	slog.DebugContext(ctx, "flags",
		"listen-addr", *listenAddr,
		"log-level", *logLevel,
		"log-format", *logFormat,
		"otel-enabled", *otelEnabled,
		"profiling-enabled", *profilingEnabled,
		"profiling-address", *profilingAddress,
		"web-framework", *webFramework,
		"statsviz-enabled", *statsvizEnabled,
	)

	switch *webFramework {
	case "echo":
		slog.DebugContext(ctx, "Starting Echo server")
		echo.Init(ctx, listenAddr, logLevel, otelEnabled, statsvizEnabled, tracer)
	case "gorilla":
		slog.DebugContext(ctx, "Starting Gorilla server")
		gorilla.Init(ctx, listenAddr, otelEnabled, statsvizEnabled, tracer)
	case "chi":
		slog.DebugContext(ctx, "Starting Chi server")
		chi.Init(ctx, listenAddr, logLevel, otelEnabled, statsvizEnabled, tracer)
	case "gin":
		slog.DebugContext(ctx, "Starting Gin server")
		gin.Init(ctx, listenAddr, logLevel, otelEnabled, statsvizEnabled, tracer, &traceProvider)
	}
}
