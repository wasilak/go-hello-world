package chi

import (
	"encoding/json"
	"net/http"

	"github.com/wasilak/go-hello-world/utils"
	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := s.SetMainResponse(ctx, r)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode main response in chi")
		appErr.AddContext("path", r.URL.Path)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) healthRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := common.HealthResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode health response in chi")
		appErr.AddContext("path", r.URL.Path)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) loggerRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	levelParam := r.URL.Query().Get("level")
	response := s.SetLogLevelResponse(ctx, levelParam)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode logger response in chi")
		appErr.AddContext("path", r.URL.Path)
		appErr.AddContext("log_level", levelParam)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) switchRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nameParam := r.URL.Query().Get("name")
	response := s.SetFrameworkResponse(ctx, nameParam)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Use the new standardized error types
		appErr := utils.WrapError(err, utils.RuntimeError, "failed to encode switch response in chi")
		appErr.AddContext("path", r.URL.Path)
		appErr.AddContext("framework", nameParam)
		appErr.LogError(ctx)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
