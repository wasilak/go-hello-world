package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"log/slog"

	otelgotracer "github.com/wasilak/otelgo/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web"
	"github.com/wasilak/go-hello-world/web/common"
	"github.com/wasilak/loggergo"
	loggergoLib "github.com/wasilak/loggergo/lib"
	loggergoTypes "github.com/wasilak/loggergo/lib/types"
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
	webFramework := flag.String("web-framework", "gorilla", "Web framework (gorilla, echo, gin, chi, fiber)")
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

	loggerConfig := loggergoTypes.Config{
		Level:        loggergoLib.LogLevelFromString(*logLevel),
		Format:       loggergoLib.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      loggergoLib.LogLevelFromString(*logLevel) == slog.LevelDebug && *logFormat == "plain",
		Output:       loggergoTypes.OutputConsole,
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
		loggerConfig.Output = loggergoTypes.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/go-hello-world"
		loggerConfig.OtelTracingEnabled = false
	}

	ctx, _, err = loggergo.Init(ctx, loggerConfig)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	logLevelConfig := loggergo.GetLogLevelAccessor()

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

	if strings.EqualFold(logLevelConfig.Level().String(), "debug") {
		slog.DebugContext(ctx, "Debug mode enabled")
	}

	frameworkOptions := common.FrameworkOptions{
		ListenAddr:      *listenAddr,
		OtelEnabled:     *otelEnabled,
		StatsvizEnabled: *statsvizEnabled,
		Tracer:          tracer,
		LogLevelConfig:  logLevelConfig,
		TraceProvider:   traceProvider,
	}

	web.RunWebServer(ctx, *webFramework, frameworkOptions)
}
