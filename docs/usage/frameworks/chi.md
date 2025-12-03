# Chi Framework Guide

This document provides specific usage information for the Chi web framework in the go-hello-world application.

## Overview

Chi is a lightweight, idiomatic and composable router for building Go HTTP services. It's built on top of the standard net/http package and provides a powerful routing API with middlewares.

## Unique Features

- **Composable Middleware**: Supports functional composition of middleware
- **URL Parameters**: Robust URL parameter extraction and routing
- **Graceful Shutdown**: Built-in support with sync.WaitGroup similar to Gorilla
- **Chi Prometheus**: Uses chi-prometheus for metrics collection
- **Standard HTTP Package**: Built on net/http for compatibility

## Performance Characteristics

- Excellent routing performance with low overhead
- Fast middleware composition
- Efficient URL parameter parsing
- Good memory usage characteristics

## Configuration Examples

### Basic Usage
```bash
go run main.go --web-framework=chi
```

### With Custom Address and Logging
```bash
go run main.go --web-framework=chi --listen-addr=0.0.0.0:8080 --log-level=WARN
```

### With Observability
```bash
go run main.go --web-framework=chi --otel-enabled=true --statsviz-enabled=true
```

### With Full Telemetry
```bash
go run main.go --web-framework=chi --otel-enabled=true --otel-host-metrics=true --otel-runtime-metrics=true --profiling-enabled=true
```

## Endpoints

All frameworks provide the same standard endpoints:

- `GET /` - Main application endpoint returning host and framework info
- `GET /health` - Health check endpoint returning "healthy" status
- `GET /logger` - Get/set logger level (query param: level)
- `GET /framework` - Get/set active framework (query param: name)
- `GET /metrics` - Prometheus metrics endpoint
- `GET /debug/statsviz/` - Statsviz visualization (when enabled)

## Best Practices

- Use Chi's middleware stack effectively
- Leverage context propagation for request-scoped values
- Monitor Chi-specific metrics for routing performance
- Use Chi's route registration patterns for organization

## Troubleshooting

### Common Issues
- Ensure proper middleware ordering in the chain
- Check sync.WaitGroup usage during shutdown
- Verify route pattern matching for complex URLs
- Monitor for context cancellation issues
