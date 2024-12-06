package chi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wasilak/go-hello-world/web"
	loggergoLib "github.com/wasilak/loggergo/lib"
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

func loggerRoute(w http.ResponseWriter, r *http.Request) {
	newLogLevelParam := r.URL.Query().Get("level")

	response := web.LoggerResponse{
		LogLevelCurrent: logLevel.Level().String(),
	}

	newLogLevel := loggergoLib.LogLevelFromString(newLogLevelParam)

	logLevel.Set(newLogLevel)

	response.LogLevelPrevious = logLevel.Level().String()

	slog.DebugContext(r.Context(), "log_level_changed", "from", response.LogLevelPrevious, "to", response.LogLevelCurrent)
	json.NewEncoder(w).Encode(response)
}
