package common

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/wasilak/go-hello-world/utils"
)

// RouteHandlerFactory provides framework-agnostic route handler patterns
type RouteHandlerFactory struct {
	WebServer *WebServer
}

// NewRouteHandlerFactory creates a new factory instance
func NewRouteHandlerFactory(ws *WebServer) *RouteHandlerFactory {
	return &RouteHandlerFactory{
		WebServer: ws,
	}
}

// MainRouteHandler returns a standardized main route handler
func (f *RouteHandlerFactory) MainRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := f.WebServer.FrameworkOptions.Tracer.Start(ctx, "mainRoute")
		defer span.End()

		response := f.WebServer.SetMainResponse(ctx, r)

		if err := sendJSONResponse(ctx, w, response, http.StatusOK); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.RuntimeError, "failed to send main route response")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)
		}
	}
}

// HealthRouteHandler returns a standardized health route handler
func (f *RouteHandlerFactory) HealthRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := f.WebServer.FrameworkOptions.Tracer.Start(ctx, "healthRoute")
		defer span.End()

		response := HealthResponse{
			Status: "healthy",
		}

		if err := sendJSONResponse(ctx, w, response, http.StatusOK); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.RuntimeError, "failed to send health route response")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)
		}
	}
}

// LoggerRouteHandler returns a standardized logger route handler
func (f *RouteHandlerFactory) LoggerRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := f.WebServer.FrameworkOptions.Tracer.Start(ctx, "loggerRoute")
		defer span.End()

		var req struct {
			Level string `json:"level"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.ValidationError, "invalid JSON in logger route request")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)

			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		response := f.WebServer.SetLogLevelResponse(ctx, strings.ToUpper(req.Level))

		if err := sendJSONResponse(ctx, w, response, http.StatusOK); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.RuntimeError, "failed to send logger route response")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)
		}
	}
}

// SwitchRouteHandler returns a standardized framework switch route handler
func (f *RouteHandlerFactory) SwitchRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := f.WebServer.FrameworkOptions.Tracer.Start(ctx, "switchRoute")
		defer span.End()

		var req struct {
			Framework string `json:"framework"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.ValidationError, "invalid JSON in switch route request")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)

			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		response := f.WebServer.SetFrameworkResponse(ctx, req.Framework)

		if err := sendJSONResponse(ctx, w, response, http.StatusOK); err != nil {
			// Use the new standardized error types
			appErr := utils.WrapError(err, utils.RuntimeError, "failed to send switch route response")
			appErr.AddContext("path", r.URL.Path)
			appErr.LogError(ctx)
		}
	}
}

// sendJSONResponse is a helper function to send JSON responses consistently
func sendJSONResponse(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.ErrorContext(ctx, "Error encoding JSON response", "error", err)
		return err
	}

	return nil
}

// GenericErrorHandler provides a standardized error handler
func (f *RouteHandlerFactory) GenericErrorHandler(err error, path string, w http.ResponseWriter) {
	// Use the new standardized error types
	appErr := utils.WrapError(err, utils.RuntimeError, "route handler error")
	appErr.AddContext("path", path)
	appErr.LogError(context.Background()) // Use background context since original may be cancelled

	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// ValidateRequest is a helper to validate incoming requests
func (f *RouteHandlerFactory) ValidateRequest(r *http.Request, expectedMethod string) *utils.AppError {
	if r.Method != expectedMethod {
		err := utils.NewAppError(
			utils.ValidationError,
			"method not allowed",
			nil,
		)
		err.AddContext("expected", expectedMethod)
		err.AddContext("received", r.Method)
		return err
	}
	return nil
}

// GetHostname returns the current hostname with error handling
func (f *RouteHandlerFactory) GetHostname() (string, *utils.AppError) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", utils.WrapError(err, utils.RuntimeError, "failed to get hostname")
	}
	return hostname, nil
}
