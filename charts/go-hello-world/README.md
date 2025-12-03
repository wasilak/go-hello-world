# go-hello-world-chart

![Version: 0.5.3](https://img.shields.io/badge/Version-0.5.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.8.2](https://img.shields.io/badge/AppVersion-1.8.2-informational?style=flat-square)

A Helm chart for go-hello-world

## Overview

go-hello-world is a reference implementation for Kubernetes deployment patterns. It includes:

- Multiple HTTP routing framework implementations (chi, echo, fiber, gin, gorilla)
- Integrated observability with OpenTelemetry
- Prometheus metrics collection
- Structured logging with slog
- Docker container support with multi-architecture builds
- Comprehensive CLI arguments for flexible configuration

This Helm chart simplifies deployment to Kubernetes and makes the application discoverable through container registries.

## Installation

### Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

### Quick Start - OCI Registry

Install directly from GitHub Container Registry (GHCR):

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
```

### Install with Custom Values

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --set replicaCount=3 \
  --set service.type=LoadBalancer
```

### Install from Local Chart

If developing locally:

```bash
helm install go-hello-world ./charts/go-hello-world
```

### Specify Version

Install a specific chart version:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart:0.5.2
```

## Configuration

The following table lists the configurable parameters of the go-hello-world chart and their default values.

### Example: Custom Replica Count

```yaml
replicaCount: 3
```

### Example: LoadBalancer Service

```yaml
service:
  type: LoadBalancer
```

### Example: Enable OpenTelemetry

```yaml
args:
  - "-listen-addr=:3000"
  - "-log-format=json"
  - "-log-level=info"
  - "-otel-enabled"

env:
  - name: "OTEL_SERVICE_NAME"
    value: "go-hello-world"
  - name: NODE_IP
    valueFrom:
      fieldRef:
        fieldPath: status.hostIP
  - name: "OTEL_EXPORTER_OTLP_ENDPOINT"
    value: "http://$(NODE_IP):4318"
```

### Example: Custom Logging Configuration

```yaml
args:
  - "-listen-addr=:3000"
  - "-log-format=json"
  - "-log-level=debug"
```

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| wasilak | <piotr.m.boruc@gmail.com> |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| annotations | object | `{}` |  |
| args[0] | string | `"-listen-addr=:3000"` |  |
| args[1] | string | `"-log-format=json"` |  |
| args[2] | string | `"-log-level=info"` |  |
| env[0].name | string | `"OTEL_SERVICE_NAME"` |  |
| env[0].value | string | `"go-hello-world"` |  |
| env[1].name | string | `"NODE_IP"` |  |
| env[1].valueFrom.fieldRef.fieldPath | string | `"status.hostIP"` |  |
| env[2].name | string | `"OTEL_EXPORTER_OTLP_ENDPOINT"` |  |
| env[2].value | string | `"http://$(NODE_IP):4318"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/wasilak/go-hello-world"` |  |
| image.tag | string | `"{{ .Chart.AppVersion }}"` |  |
| imagePullPolicy | string | `"Always"` |  |
| podAnnotations | object | `{}` |  |
| replicaCount | string | `"1"` |  |
| service.name | string | `"go-hello-world"` |  |
| service.port | int | `3000` |  |
| service.targetPort | int | `3000` |  |
| service.type | string | `"ClusterIP"` |  |

## Accessing the Application

After installation, you can access the application in several ways:

### Port Forward (for testing)

```bash
kubectl port-forward svc/go-hello-world 3000:3000
```

Then open http://localhost:3000 in your browser.

### Via Ingress

Create an ingress resource to expose the application:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-hello-world
spec:
  ingressClassName: nginx
  rules:
    - host: go-hello-world.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: go-hello-world
                port:
                  number: 3000
```

### Via LoadBalancer

If you set `service.type=LoadBalancer`:

```bash
kubectl get svc go-hello-world
```

Use the external IP to access the application.

## Application Arguments

The application supports the following command-line arguments:

| Argument | Default | Description |
|----------|---------|-------------|
| `-listen-addr` | `:3000` | Address to listen on |
| `-log-format` | `json` | Log format (json, text) |
| `-log-level` | `info` | Log level (debug, info, warn, error) |
| `-otel-enabled` | `false` | Enable OpenTelemetry instrumentation |

## Upgrading

To upgrade an existing installation:

```bash
helm upgrade go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
```

With custom values:

```bash
helm upgrade go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --set replicaCount=5
```

## Uninstalling

To remove the release:

```bash
helm uninstall go-hello-world
```

## Troubleshooting

### Pod is not starting

Check pod status and logs:

```bash
# Check pod status
kubectl get pods -l app=go-hello-world

# View pod logs
kubectl logs -l app=go-hello-world

# Describe pod for events
kubectl describe pod -l app=go-hello-world
```

### Application is not accessible

Verify service is running:

```bash
kubectl get svc go-hello-world
kubectl get endpoints go-hello-world
```

Test connectivity:

```bash
kubectl run -it --rm debug --image=busybox --restart=Never -- wget -O- http://go-hello-world:3000
```

### Image pull errors

Verify image repository and tag:

```bash
# Check image pull policy
kubectl get deployment go-hello-world -o yaml | grep -A5 image

# Pull the image manually to test
docker pull ghcr.io/wasilak/go-hello-world:1.8.2
```

## Support

- 📖 [Project Repository](https://github.com/wasilak/go-hello-world)
- 🐛 [Issue Tracker](https://github.com/wasilak/go-hello-world/issues)
- 📦 [OCI Registry](oci://ghcr.io/wasilak/go-hello-world-chart)

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
