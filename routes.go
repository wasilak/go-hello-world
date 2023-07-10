package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "ok"}
	_, spanJsonEncode := tracer.Start(r.Context(), "json encode response")
	json.NewEncoder(w).Encode(response)
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
	slog.DebugCtx(ctx, "rootHandler", "response", response)

	spanSession.End()

	_, spanJsonEncode := tracer.Start(ctx, "json encode response")
	json.NewEncoder(w).Encode(response)
	spanJsonEncode.End()
}
