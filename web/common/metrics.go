package common

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterCollectorIfNotRegistered registers a collector only if it's not already registered.
// This prevents duplicate registration errors when multiple packages try to register
// the same collectors (e.g., Go runtime metrics, process metrics).
func RegisterCollectorIfNotRegistered(c prometheus.Collector) {
	// Try to register the collector
	err := prometheus.Register(c)
	if err != nil {
		// If registration fails due to duplicate, it's already registered - this is fine
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return
		}
		// For other errors, panic as this indicates a real problem
		panic(err)
	}
	// Registration succeeded, but we need to unregister it since we were just checking
	prometheus.Unregister(c)
}
