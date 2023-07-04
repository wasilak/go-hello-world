package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	InitTracer(ctx)
	// if otelEnabled {
	// }

	ctx, newSpan := tracer.Start(ctx, "rootHandler")
	// defer newSpan.End()

	session, err := store.Get(r, "session-go-hello-world")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	APIStatsFromSession := session.Values["apistats"]

	var ok bool
	var response APIResponse

	response.APIStats, ok = APIStatsFromSession.(APIStats)

	if !ok {
		log.Println("session not initialized (yet)")
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

	session.Values["apistats"] = response.APIStats

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.InfoCtx(ctx, "root")
	newSpan.End()

	json.NewEncoder(w).Encode(response)
}
