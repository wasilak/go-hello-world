package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ErrorType represents the category of an error
type ErrorType string

const (
	// ConfigError represents configuration-related errors
	ConfigError ErrorType = "config"

	// RuntimeError represents runtime errors during application execution
	RuntimeError ErrorType = "runtime"

	// FrameworkError represents errors specific to web frameworks
	FrameworkError ErrorType = "framework"

	// ValidationError represents validation errors
	ValidationError ErrorType = "validation"
)

// AppError represents a standardized application error with type and context
type AppError struct {
	Type       ErrorType
	Message    string
	Err        error
	OccurredAt time.Time
	Context    map[string]interface{}
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError instance
func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		Err:        err,
		OccurredAt: time.Now(),
		Context:    make(map[string]interface{}),
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, errType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		Err:        err,
		OccurredAt: time.Now(),
		Context:    make(map[string]interface{}),
	}
}

// AddContext adds context information to an AppError
func (e *AppError) AddContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// ErrorCounterVec is a Prometheus counter vector for tracking errors by type
var ErrorCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_errors_total",
		Help: "Total number of application errors by type",
	},
	[]string{"error_type", "error_message"},
)

// LogError logs the AppError using slog and updates metrics
func (e *AppError) LogError(ctx interface{}) {
	// Update error metrics
	ErrorCounterVec.WithLabelValues(string(e.Type), e.Message).Inc()

	// Log error with context
	loggerAttrs := []interface{}{
		"error_type", e.Type,
		"message", e.Message,
		"occurred_at", e.OccurredAt,
	}

	// Add context fields to log
	for key, value := range e.Context {
		loggerAttrs = append(loggerAttrs, key, value)
	}

	if e.Err != nil {
		loggerAttrs = append(loggerAttrs, "original_error", e.Err)
	}

	// Use slog to log the error with context
	slog.ErrorContext(ctx.(context.Context), e.Error(), loggerAttrs...)
}

// IsConfigError checks if an error is of configuration type
func IsConfigError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ConfigError
	}
	return false
}

// IsRuntimeError checks if an error is of runtime type
func IsRuntimeError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == RuntimeError
	}
	return false
}

// IsFrameworkError checks if an error is of framework type
func IsFrameworkError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == FrameworkError
	}
	return false
}

// IsValidationError checks if an error is of validation type
func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ValidationError
	}
	return false
}

// AsAppError tries to convert an error to AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
