package gorilla

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/wasilak/go-hello-world/utils"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: fmt.Sprintf("%s_http_duration_seconds", strings.ReplaceAll(utils.GetAppName(), "-", "_")),
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_requests_count_total", strings.ReplaceAll(utils.GetAppName(), "-", "_")),
		Help: "HTTP requests count.",
	}, []string{"path", "host"})
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

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

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		requestCounter.With(prometheus.Labels{"path": path, "host": r.Host}).Inc()
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}
