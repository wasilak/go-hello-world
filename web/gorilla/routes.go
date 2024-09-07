package gorilla

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"log/slog"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "healthHandler called")
	w.WriteHeader(http.StatusOK)
	response := HealthResponse{Status: "ok"}
	_, spanJsonEncode := tracer.Start(r.Context(), "json encode response")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
	}
	spanJsonEncode.End()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	ctx, spanSession := tracer.Start(r.Context(), "session")

	ctx, spanResponse := tracer.Start(ctx, "response")
	var response APIResponse

	hostname, _ := os.Hostname()
	response.Host = hostname

	response.Request = APIResponseRequest{
		Host:       r.Host,
		URL:        r.URL,
		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
		Method:     r.Method,
		Proto:      r.Proto,
		UserAgent:  r.UserAgent(),
		Headers:    r.Header,
	}
	spanResponse.End()

	spanSession.AddEvent(fmt.Sprintf("%+v", response))
	slog.DebugContext(ctx, "rootHandler", "response", response)

	spanSession.End()

	_, spanJsonEncode := tracer.Start(ctx, "json encode response")
	json.NewEncoder(w).Encode(response)
	spanJsonEncode.End()
}
