# Configuration Flags

This document provides a comprehensive reference for all command-line flags available in the go-hello-world application.

## Available Flags

### Server Configuration

#### `--listen-addr`
- **Type**: String
- **Default**: `127.0.0.1:3000`
- **Description**: Server listen address
- **Example**: `--listen-addr=0.0.0.0:8080`

### Logging Configuration

#### `--log-level`
- **Type**: String
- **Default**: `INFO`
- **Description**: Log level (options: DEBUG, INFO, WARN, ERROR)
- **Example**: `--log-level=DEBUG`

#### `--log-format`
- **Type**: String
- **Default**: `text`
- **Description**: Log format (options: text, json, plain)
- **Example**: `--log-format=json`

#### `--dev-flavor`
- **Type**: String
- **Default**: `tint`
- **Description**: Development flavor for logging (options: tint, log/slog)
- **Example**: `--dev-flavor=log/slog`

#### `--output-type`
- **Type**: String
- **Default**: `console`
- **Description**: Output type for logging (options: console, file, fanout)
- **Example**: `--output-type=file`

### Observability Configuration

#### `--otel-enabled`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable OpenTelemetry traces
- **Example**: `--otel-enabled=true`

#### `--otel-host-metrics`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable OpenTelemetry host metrics
- **Example**: `--otel-host-metrics=true`

#### `--otel-runtime-metrics`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable OpenTelemetry runtime metrics
- **Example**: `--otel-runtime-metrics=true`

### Performance and Profiling

#### `--statsviz-enabled`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable statsviz for visualization
- **Example**: `--statsviz-enabled=true`

#### `--profiling-enabled`
- **Type**: Boolean
- **Default**: `false`
- **Description**: Enable profiling
- **Example**: `--profiling-enabled=true`

#### `--profiling-address`
- **Type**: String
- **Default**: `127.0.0.1:4040`
- **Description**: Profiling server address
- **Example**: `--profiling-address=0.0.0.0:4040`

### Framework Configuration

#### `--web-framework`
- **Type**: String
- **Default**: `gorilla`
- **Description**: Web framework to use (options: gorilla, echo, gin, chi, fiber)
- **Example**: `--web-framework=echo`

## Usage Examples

### Basic Usage
```bash
go run main.go
```

### With Custom Address and Debug Logging
```bash
go run main.go --listen-addr=0.0.0.0:8080 --log-level=DEBUG
```

### With OpenTelemetry and Profiling
```bash
go run main.go --otel-enabled=true --profiling-enabled=true --profiling-address=0.0.0.0:4040
```

### With Different Web Framework
```bash
go run main.go --web-framework=gin
```

## Validation Rules

- `--listen-addr` must be a valid host:port combination
- `--log-level` values are case-insensitive and validated against supported levels
- `--log-format` must be one of the supported formats
- `--profiling-address` must be a valid host:port combination
- `--web-framework` must be one of the supported web frameworks