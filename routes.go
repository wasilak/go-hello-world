package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/codes"
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

	session, err := store.Get(r, "session-go-hello-world")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		spanSession.SetStatus(codes.Error, "store.Get() failed")
		spanSession.RecordError(err)
		return
	}

	APIStatsFromSession := session.Values["apistats"]

	ctx, spanResponse := tracer.Start(ctx, "response")
	var ok bool
	var response APIResponse

	response.APIStats, ok = APIStatsFromSession.(APIStats)

	if !ok {
		slog.DebugCtx(ctx, "session not initialized (yet)")
	}

	response.APIStats.Counter++

	hostname, _ := os.Hostname()
	response.Host = hostname

	if nil == response.APIStats.Hostnames {
		response.APIStats.Hostnames = make(map[string]int)
	}

	response.APIStats.Hostnames[hostname]++

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

	session.Values["apistats"] = response.APIStats

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spanSession.AddEvent(fmt.Sprintf("%+v", response))
	slog.DebugCtx(ctx, "rootHandler", "response", response)

	spanSession.End()

	_, spanJsonEncode := tracer.Start(ctx, "json encode response")
	json.NewEncoder(w).Encode(response)
	spanJsonEncode.End()
}
