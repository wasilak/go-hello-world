package gorilla

import (
	"encoding/json"
	"net/http"
	"os"

	"log/slog"

	"github.com/wasilak/go-hello-world/web"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "healthHandler called")
	w.WriteHeader(http.StatusOK)
	response := web.HealthResponse{Status: "ok"}
	_, spanJsonEncode := tracer.Start(r.Context(), "json encode response")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
	}
	spanJsonEncode.End()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	ctx, spanResponse := tracer.Start(r.Context(), "response")

	hostname, _ := os.Hostname()
	response := web.APIResponse{
		Host: hostname,
		Request: web.APIResponseRequest{
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
	spanResponse.End()

	slog.DebugContext(ctx, "rootHandler", "response", response)

	_, spanJsonEncode := tracer.Start(ctx, "json encode response")
	json.NewEncoder(w).Encode(response)
	spanJsonEncode.End()
}
