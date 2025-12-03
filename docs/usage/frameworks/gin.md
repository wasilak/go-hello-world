# Gin Framework Guide

This document provides specific usage information for the Gin web framework in the go-hello-world application.

## Overview

Gin is a HTTP web framework written in Go that features a Martini-like API with much better performance. The implementation uses Gin's fast routing and built-in middleware capabilities.

## Unique Features

- **Exceptional Performance**: One of the fastest Go frameworks with minimal overhead
- **Gin Promtheus**: Integration with gin-prometheus for metrics collection
- **Built-in Logger**: Custom logging with support for different flavors
- **Recovery Middleware**: Automatic panic recovery with logging

## Performance Characteristics

- Outstanding request handling performance
- Fast route matching and execution
- Low memory allocation per request
- Optimized for high-concurrency scenarios

## Configuration Examples

### Basic Usage
```bash
go run main.go --web-framework=gin
```

### With Debug Mode
```bash
go run main.go --web-framework=gin --log-level=DEBUG --dev-flavor=log/slog
```

### With Observability
```bash
go run main.go --web-framework=gin --otel-enabled=true --otel-host-metrics=true
```

### With Custom Settings
```bash
go run main.go --web-framework=gin --listen-addr=0.0.0.0:8080 --log-format=json --output-type=console
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

- Use Gin's built-in validator for request validation
- Leverage Gin's HTML templating system for server-side rendering
- Monitor Gin-specific metrics for performance insights
- Use Gin's error management system for consistent error handling

## Troubleshooting

### Common Issues
- Ensure gin.DefaultWriter is properly configured for logging
- Check Gin mode (debug vs release) settings
- Verify middleware ordering for proper request processing
- Monitor for memory issues with large request bodies
