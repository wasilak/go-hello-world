# Echo Framework Guide

This document provides specific usage information for the Echo web framework in the go-hello-world application.

## Overview

Echo is a high-performance, minimalist Go web framework that provides a robust set of features for building APIs and web applications. The implementation uses Echo's middleware system and performance-optimized routing.

## Unique Features

- **High Performance**: Optimized for speed and low memory usage
- **Middleware System**: Built-in middleware for Gzip, CORS, logging, and recovery
- **Echo Prometheus**: Uses echo-contrib/echoprometheus for metrics collection
- **OpenTelemetry Support**: Full tracing integration with otelecho middleware

## Performance Characteristics

- Excellent performance with low memory footprint
- Fast routing with optimized middleware pipeline
- Efficient JSON serialization/deserialization
- Minimal overhead for standard operations

## Configuration Examples

### Basic Usage
```bash
go run main.go --web-framework=echo
```

### With Performance Optimization
```bash
go run main.go --web-framework=echo --listen-addr=0.0.0.0:8080 --log-level=INFO
```

### With Observability
```bash
go run main.go --web-framework=echo --otel-enabled=true --otel-runtime-metrics=true --log-level=DEBUG
```

### With Profiling and Statsviz
```bash
go run main.go --web-framework=echo --profiling-enabled=true --statsviz-enabled=true
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

- Use Echo's validator for request validation
- Leverage Echo's centralized error handling
- Monitor echo-specific metrics for performance insights
- Use Echo's grouping feature for organizing route handlers

## Troubleshooting

### Common Issues
- Ensure echo.WrapHandler is properly used for non-echo handlers
- Check middleware order in the pipeline
- Verify Echo's context propagation in middleware chain
- Monitor memory usage with high concurrency
