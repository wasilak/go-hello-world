# Gorilla Framework Guide

This document provides specific usage information for the Gorilla web framework in the go-hello-world application.

## Overview

Gorilla is a toolkit of packages for writing web applications in Go. The implementation uses gorilla/mux for routing with robust pattern matching and request handling.

## Unique Features

- **Robust Routing**: Uses gorilla/mux router with advanced pattern matching
- **Graceful Shutdown**: Implements graceful shutdown using sync.WaitGroup for proper request handling during shutdown
- **Middleware Pipeline**: Custom Prometheus middleware with request duration and count metrics
- **OpenTelemetry Integration**: Full tracing support with otelmux middleware

## Performance Characteristics

- Good performance with established routing
- sync.WaitGroup ensures all requests complete during shutdown
- Custom Prometheus metrics collection
- Memory efficient for typical request loads

## Configuration Examples

### Basic Usage
```bash
go run main.go --web-framework=gorilla
```

### With Custom Address
```bash
go run main.go --web-framework=gorilla --listen-addr=0.0.0.0:8080
```

### With Observability
```bash
go run main.go --web-framework=gorilla --otel-enabled=true --log-level=DEBUG
```

### With Profiling
```bash
go run main.go --web-framework=gorilla --profiling-enabled=true --profiling-address=0.0.0.0:4040
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

- Use gorilla/mux for complex route patterns
- Leverage sync.WaitGroup for graceful shutdown during deployment
- Monitor custom Prometheus metrics for performance tuning
- Enable OpenTelemetry for distributed tracing in microservices

## Troubleshooting

### Common Issues
- Ensure sync.WaitGroup is properly managed during shutdown
- Check Prometheus metrics for routing performance
- Verify middleware order in the chain
