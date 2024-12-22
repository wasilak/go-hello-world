package common

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"

	loggergoLib "github.com/wasilak/loggergo/lib"
	"go.opentelemetry.io/otel/trace"
)

// HealthResponse type
type HealthResponse struct {
	Status string `json:"status"`
}

// LoggerResponse type
type LoggerResponse struct {
	LogLevelCurrent  string `json:"log_level_current"`
	LogLevelPrevious string `json:"log_level_previous"`
}

// FrameworkResponse type
type FrameworkResponse struct {
	FrameworkCurrent  string `json:"framework_current"`
	FrameworkPrevious string `json:"framework_previous"`
}

// APIResponseRequest type
type APIResponseRequest struct {
	Host       string      `json:"host"`
	RemoteAddr string      `json:"remote_addr"`
	RequestURI string      `json:"request_uri"`
	Method     string      `json:"method"`
	Proto      string      `json:"proto"`
	UserAgent  string      `json:"user_agent"`
	URL        *url.URL    `json:"url"`
	Headers    http.Header `json:"headers"`
}

// APIResponse type
type APIResponse struct {
	Host      string             `json:"host"`
	Framework string             `json:"framework"`
	Request   APIResponseRequest `json:"request"`
}

type FrameworkOptions struct {
	ListenAddr      string
	OtelEnabled     bool
	StatsvizEnabled bool
	Tracer          trace.Tracer
	LogLevelConfig  *slog.LevelVar
	TraceProvider   trace.TracerProvider
}

type WebServer struct {
	MU               sync.Mutex
	Running          bool
	Framework        string
	FrameworkOptions FrameworkOptions
}

type WebServerInterface interface {
	Start(context.Context)
	Stop(ctx context.Context)
}

// Create a channel to signal framework changes
var (
	FrameworkChannel chan string
)

func (w *WebServer) SetMainResponse(ctx context.Context, r *http.Request) APIResponse {
	_, span := w.FrameworkOptions.Tracer.Start(ctx, "response")
	hostname, _ := os.Hostname()
	response := APIResponse{
		Host:      hostname,
		Framework: w.Framework,
		Request: APIResponseRequest{
			Host:       r.Host,
			URL:        r.URL,
			RemoteAddr: r.RemoteAddr,
			RequestURI: r.RequestURI,
			Method:     r.Method,
			Proto:      r.Proto,
			UserAgent:  r.UserAgent(),
			Headers:    r.Header,
		},
	}
	span.End()
	return response
}

func (w *WebServer) SetLogLevelResponse(ctx context.Context, current string) LoggerResponse {
	ctx, span := w.FrameworkOptions.Tracer.Start(ctx, "logLevelResponse")

	response := LoggerResponse{
		LogLevelCurrent: w.FrameworkOptions.LogLevelConfig.Level().String(),
	}

	newLogLevel := loggergoLib.LogLevelFromString(current)

	w.FrameworkOptions.LogLevelConfig.Set(newLogLevel)

	response.LogLevelPrevious = w.FrameworkOptions.LogLevelConfig.Level().String()

	slog.DebugContext(ctx, "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)

	span.End()

	return response
}

func (w *WebServer) SetFrameworkResponse(ctx context.Context, current string) FrameworkResponse {
	ctx, span := w.FrameworkOptions.Tracer.Start(ctx, "FrameworkResponse")

	response := FrameworkResponse{
		FrameworkPrevious: w.Framework,
	}

	if w.Framework != current {
		FrameworkChannel <- current
	}

	response.FrameworkCurrent = current

	slog.DebugContext(ctx, "framework_not_changed", "from", response.FrameworkPrevious, "to", response.FrameworkCurrent)

	span.End()

	return response
}
