package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"log/slog"

	otelgotracer "github.com/wasilak/otelgo/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web"
	"github.com/wasilak/go-hello-world/web/common"
	"github.com/wasilak/loggergo"
	"github.com/wasilak/profilego"
	"github.com/wasilak/profilego/config"
	"github.com/wasilak/profilego/core"
)

func main() {
	// Create a context that listens for system signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		sig := <-sigChan
		slog.DebugContext(ctx, "Received signal, shutting down", "signal", sig.String())
		cancel() // Cancel the context
	}()

	listenAddr := flag.String("listen-addr", "127.0.0.1:3000", "server listen address")
	logLevel := flag.String("log-level", slog.LevelInfo.String(), fmt.Sprintf("log level %s", loggergo.Types.AllLogLevels()))
	logFormat := flag.String("log-format", loggergo.Types.LogFormatText.String(), fmt.Sprintf("log format %s", loggergo.Types.AllLogFormats()))
	otelEnabled := flag.Bool("otel-enabled", false, "OpenTelemetry traces enabled")
	otelHostMetricsEnabled := flag.Bool("otel-host-metrics", false, "OpenTelemetry host metrics enabled")
	otelRuntimeMetricsEnabled := flag.Bool("otel-runtime-metrics", false, "OpenTelemetry runtime metrics enabled")
	statsvizEnabled := flag.Bool("statsviz-enabled", false, "statsviz enabled")
	profilingEnabled := flag.Bool("profiling-enabled", false, "Profiling enabled")
	profilingAddress := flag.String("profiling-address", "http://localhost:4040", "Profiling address")
	webFramework := flag.String("web-framework", "gorilla", "Web framework (gorilla, echo, gin, chi, fiber)")
	devFlavor := flag.String("dev-flavor", loggergo.Types.DevFlavorTint.String(), fmt.Sprintf("Dev flavor %s", loggergo.Types.AllDevFlavors()))
	outPutType := flag.String("output-type", loggergo.Types.OutputConsole.String(), fmt.Sprintf("Output type %s", loggergo.Types.AllOutputTypes()))
	flag.Parse()

	if *profilingEnabled {
		profileGoConfig := config.Config{
			ApplicationName: utils.GetAppName(),
			ServerAddress:   *profilingAddress,
			Backend:         core.PyroscopeBackend,
			InitialState:    core.ProfilingEnabled,
		}

		err := profilego.InitWithConfig(profileGoConfig)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to initialize profiling: %v", "error", err)
			os.Exit(1)
		}
	}

	loggerConfig := loggergo.Config{
		Level:        loggergo.Types.LogLevelFromString(*logLevel),
		Format:       loggergo.Types.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      loggergo.Types.LogLevelFromString(*logLevel) == slog.LevelDebug && *logFormat == "plain",
		Output:       loggergo.Types.OutputTypeFromString(*outPutType),
		DevFlavor:    loggergo.Types.DevFlavorFromString(*devFlavor),
	}

	var tracer = otel.Tracer(utils.GetAppName())
	var traceProvider *trace.TracerProvider
	var err error

	if *otelEnabled {
		otelGoTracingConfig := otelgotracer.Config{
			HostMetricsEnabled:    *otelHostMetricsEnabled,
			RuntimeMetricsEnabled: *otelRuntimeMetricsEnabled,
		}
		ctx, traceProvider, err = otelgotracer.Init(ctx, otelGoTracingConfig)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			os.Exit(1)
		}

		loggerConfig.OtelServiceName = utils.GetAppName()
		loggerConfig.Output = loggergo.Types.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/go-hello-world"
		loggerConfig.OtelTracingEnabled = true

		otel.SetTracerProvider(traceProvider)
		tracer = traceProvider.Tracer(utils.GetAppName())

		defer func() {
			if err := traceProvider.Shutdown(ctx); err != nil {
				slog.ErrorContext(ctx, "Failed to shut down trace provider", "error", err)
			}
		}()
	}

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

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

	// Create a channel to signal framework changes
	common.FrameworkChannel = make(chan string)

	go web.RunWebServer(ctx, frameworkOptions)

	common.FrameworkChannel <- *webFramework

	// Wait for the context to be canceled
	<-ctx.Done()

	// Perform any necessary cleanup here
	slog.InfoContext(ctx, "Application exiting")
}
