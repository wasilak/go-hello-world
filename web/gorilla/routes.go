package gorilla

import (
	"encoding/json"
	"net/http"

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
	response := web.ConstructResponse(r)
	spanResponse.End()

	slog.DebugContext(ctx, "rootHandler", "response", response)

	_, spanJsonEncode := tracer.Start(ctx, "json encode response")
	json.NewEncoder(w).Encode(response)
	spanJsonEncode.End()
}
