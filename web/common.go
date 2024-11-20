package web

import (
	"net/http"
	"net/url"
	"os"
)

// HealthResponse type
type HealthResponse struct {
	Status string `json:"status"`
}

// LoggerResponse type
type LoggerResponse struct {
	LogLevelCurrent  string `json:"log_level_current"`
	LogLevelPrevious string `json:"log_level_previous"`
}

// APIResponseRequest type
type APIResponseRequest struct {
	Host       string      `json:"host"`
	RemoteAddr string      `json:"remote_addr"`
	RequestURI string      `json:"request_uri"`
	Method     string      `json:"method"`
	Proto      string      `json:"proto"`
	UserAgent  string      `json:"user_agent"`
	URL        *url.URL    `json:"url"`
	Headers    http.Header `json:"headers"`
}

// APIResponse type
type APIResponse struct {
	Host    string             `json:"host"`
	Request APIResponseRequest `json:"request"`
}

func ConstructResponse(r *http.Request) APIResponse {
	hostname, _ := os.Hostname()
	response := APIResponse{
		Host: hostname,
		Request: APIResponseRequest{
			Host:       r.Host,
			URL:        r.URL,
			RemoteAddr: r.RemoteAddr,
			RequestURI: r.RequestURI,
			Method:     r.Method,
			Proto:      r.Proto,
			UserAgent:  r.UserAgent(),
			Headers:    r.Header,
		},
	}
	return response
}
