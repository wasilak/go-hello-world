# Fiber Framework Guide

This document provides specific usage information for the Fiber web framework in the go-hello-world application.

## Overview

Fiber is an Express.js inspired web framework built on Fasthttp, the fastest HTTP engine for Go. It's designed to be intuitive, fast, and compatible with existing middleware.

## Unique Features

- **Based on Fasthttp**: Faster than standard net/http, inspired by Express.js
- **Middleware Support**: Full middleware ecosystem with built-in compression
- **Fiber Prometheus**: Direct Prometheus metrics integration
- **Express-like Syntax**: Familiar API for developers from Node.js background

## Performance Characteristics

- Exceptional performance due to Fasthttp engine
- Very low memory allocation
- High throughput with concurrent requests
- Optimized for microservice architectures

## Configuration Examples

### Basic Usage
```bash
go run main.go --web-framework=fiber
```

### With Performance Focus
```bash
go run main.go --web-framework=fiber --log-level=INFO --output-type=console
```

### With Observability
```bash
go run main.go --web-framework=fiber --otel-enabled=true --profiling-enabled=true --statsviz-enabled=true
```

### With Full Configuration
```bash
go run main.go --web-framework=fiber --listen-addr=0.0.0.0:8080 --log-level=DEBUG --log-format=json --dev-flavor=tint
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

- Leverage Fiber's performance optimizations
- Use Fiber's built-in validation and error handling
- Monitor Fasthttp-specific metrics
- Use Fiber's grouping and mounting features for organization

## Troubleshooting

### Common Issues
- Be aware of Fasthttp's different behavior than net/http
- Check fiber context lifecycle in middleware
- Verify adaptor/v2 usage for non-Fiber handlers
- Monitor goroutine usage with high concurrency
