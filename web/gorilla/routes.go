package gorilla

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/wasilak/go-hello-world/web"
	"github.com/wasilak/loggergo"
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

func loggerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the new log level parameter from the query
	newLogLevelParam := r.URL.Query().Get("level")

	// Prepare the response structure
	response := web.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	// Parse and set the new log level
	newLogLevel := loggergo.LogLevelFromString(newLogLevelParam)
	logLevel.Set(newLogLevel)

	// Update the response with the previous log level (after setting the new one)
	response.LogLevelPrevious = logLevel.Level().String()

	// Log the change using structured logging
	slog.DebugContext(r.Context(), "log_level_changed",
		"from", response.LogLevelPrevious,
		"to", response.LogLevelCurrent,
	)

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log encoding failure
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
