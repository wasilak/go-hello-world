package chi

import (
	"encoding/json"
	"net/http"

	"github.com/wasilak/go-hello-world/web/common"
)

func (s *Server) mainRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.SetMainResponse(r.Context(), r))
}

func (s *Server) healthRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(common.HealthResponse{Status: "ok"})
}

func (s *Server) loggerRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.SetLogLevelResponse(r.Context(), r.URL.Query().Get("level")))
}

func (s *Server) switchRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.SetFrameworkResponse(r.Context(), r.URL.Query().Get("name")))
}
