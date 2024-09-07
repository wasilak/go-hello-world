package web

import (
	"net/http"
	"net/url"
)

// HealthResponse type
type HealthResponse struct {
	Status string `json:"status"`
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
