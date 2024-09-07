package chi

import (
	"encoding/json"
	"net/http"

	"github.com/wasilak/go-hello-world/web"
)

func mainRoute(w http.ResponseWriter, r *http.Request) {
	_, spanResponse := tracer.Start(r.Context(), "response")
	response := web.ConstructResponse(r)
	spanResponse.End()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthRoute(w http.ResponseWriter, r *http.Request) {
	response := web.HealthResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
