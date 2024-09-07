package gorilla

import (
	"net/http"
	"net/url"
	"time"

	"log/slog"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

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

// Chain applies middlewares to a http.HandlerFunc
func chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Logging logs all requests with its path and the time it took to process
func logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() {
				slog.InfoContext(r.Context(), "request", "path", r.URL.Path, "time", time.Since(start))
			}()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
