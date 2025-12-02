package gorilla

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "healthHandler called")
	w.WriteHeader(http.StatusOK)
	response := common.HealthResponse{Status: "ok"}

	// Ensure context is used by making sure FrameworkOptions is accessible
	if s.FrameworkOptions.Tracer != nil {
		_, span := s.FrameworkOptions.Tracer.Start(ctx, "json encode response")
		defer span.End()
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode health response")
		appErr.AddContext("path", r.URL.Path)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := s.SetMainResponse(ctx, r)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode main response")
		appErr.AddContext("path", r.URL.Path)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) loggerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Extract the new log level parameter from the query
	newLogLevelParam := r.URL.Query().Get("level")

	response := s.SetLogLevelResponse(ctx, newLogLevelParam)

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode logger response")
		appErr.AddContext("path", r.URL.Path)
		appErr.AddContext("log_level", newLogLevelParam)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) switchRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Extract the new log level parameter from the query
	newFrameworkParam := r.URL.Query().Get("name")

	response := s.SetFrameworkResponse(ctx, newFrameworkParam)

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode switch response")
		appErr.AddContext("path", r.URL.Path)
		appErr.AddContext("framework", newFrameworkParam)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
