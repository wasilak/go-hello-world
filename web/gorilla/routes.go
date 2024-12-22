package gorilla

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	slog.DebugContext(r.Context(), "healthHandler called")
	w.WriteHeader(http.StatusOK)
	response := common.HealthResponse{Status: "ok"}
	_, spanJsonEncode := s.FrameworkOptions.Tracer.Start(r.Context(), "json encode response")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
	}
	spanJsonEncode.End()
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.SetMainResponse(r.Context(), r))
}

func (s *Server) loggerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the new log level parameter from the query
	newLogLevelParam := r.URL.Query().Get("level")

	response := s.SetLogLevelResponse(r.Context(), newLogLevelParam)

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log encoding failure
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) switchRoute(w http.ResponseWriter, r *http.Request) {
	// Extract the new log level parameter from the query
	newFrameworkParam := r.URL.Query().Get("name")

	response := s.SetFrameworkResponse(r.Context(), newFrameworkParam)

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log encoding failure
		slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
