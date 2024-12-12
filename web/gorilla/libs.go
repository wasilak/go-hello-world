package gorilla

import (
	"fmt"
	"net/http"
	"strings"

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
