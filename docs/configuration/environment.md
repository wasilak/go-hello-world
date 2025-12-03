# Environment Variables

This document provides a comprehensive reference for all environment variables available in the go-hello-world application.

## Available Environment Variables

### Application Configuration

#### `APP_NAME`
- **Type**: String
- **Default**: `go-hello-world`
- **Description**: Name of the application, used for logging and tracing
- **Usage**: Provides fallback when `OTEL_SERVICE_NAME` is not set
- **Example**: `APP_NAME=my-go-app`

#### `OTEL_SERVICE_NAME`
- **Type**: String
- **Default**: None (falls back to `APP_NAME` or `go-hello-world`)
- **Description**: Service name for OpenTelemetry tracing
- **Usage**: Takes precedence over `APP_NAME` when set
- **Example**: `OTEL_SERVICE_NAME=my-service`

## Hierarchy and Priority

The application follows this priority order for determining the application name:
1. `OTEL_SERVICE_NAME` (highest priority)
2. `APP_NAME` (medium priority)
3. `go-hello-world` (default, when neither environment variable is set)

## Relationship to Command Line Flags

Environment variables provide alternative configuration to command line flags:

| Environment Variable | Related Flag | Description |
|----------------------|--------------|-------------|
| `OTEL_SERVICE_NAME` / `APP_NAME` | N/A (used internally) | Application name used for service identification |

## Usage Examples

### Using Environment Variables
```bash
export APP_NAME=my-go-hello-world
go run main.go
```

### With Docker
```bash
docker run -e APP_NAME=my-app -e OTEL_SERVICE_NAME=my-service my-go-hello-world
```

### In Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-hello-world
spec:
  template:
    spec:
      containers:
      - name: app
        image: go-hello-world
        env:
        - name: APP_NAME
          value: "production-go-app"
        - name: OTEL_SERVICE_NAME
          value: "my-production-service"
```

## Validation Rules

- `APP_NAME` and `OTEL_SERVICE_NAME` accept any non-empty string value
- If both are set, `OTEL_SERVICE_NAME` takes precedence
- Empty values are treated as unset, triggering fallback behavior
- Values containing special characters should be properly quoted when used in shell scripts
