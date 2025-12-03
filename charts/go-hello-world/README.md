# go-hello-world Helm Chart

![Version: 0.5.2](https://img.shields.io/badge/Version-0.5.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.8.2](https://img.shields.io/badge/AppVersion-1.8.2-informational?style=flat-square)

A Kubernetes Helm chart for deploying go-hello-world, a simple web application for testing container orchestrators, autoscalers, and infrastructure components.

## Overview

go-hello-world is a reference implementation for Kubernetes deployment patterns. It includes:
- Multiple HTTP routing framework implementations (chi, echo, fiber, gin, gorilla)
- Integrated observability with OpenTelemetry
- Prometheus metrics collection
- Structured logging with slog
- Docker container support with multi-architecture builds

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
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --version 0.5.2
```

## Configuration

The following table lists configurable parameters and their default values.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `replicaCount` | string | `"1"` | Number of pod replicas |
| `image.repository` | string | `"ghcr.io/wasilak/go-hello-world"` | Container image repository |
| `image.pullPolicy` | string | `"IfNotPresent"` | Container image pull policy |
| `image.tag` | string | `"{{ .Chart.AppVersion }}"` | Container image tag (defaults to chart appVersion) |
| `imagePullPolicy` | string | `"Always"` | Global image pull policy |
| `service.name` | string | `"go-hello-world"` | Kubernetes service name |
| `service.type` | string | `"ClusterIP"` | Service type (`ClusterIP`, `NodePort`, `LoadBalancer`) |
| `service.port` | number | `3000` | Service port |
| `service.targetPort` | number | `3000` | Container target port |
| `args` | array | `["-listen-addr=:3000", "-log-format=json", "-log-level=info"]` | Application command-line arguments |
| `env` | array | `[{OTEL_SERVICE_NAME, NODE_IP, OTEL_EXPORTER_OTLP_ENDPOINT}]` | Environment variables for pods |
| `annotations` | object | `{}` | Pod annotations |
| `podAnnotations` | object | `{}` | Additional pod annotations |

### Example: Custom Configuration

Create a `values-custom.yaml`:

```yaml
replicaCount: 3

image:
  tag: "1.8.2"

service:
  type: LoadBalancer
  port: 8080
  targetPort: 3000

args:
  - "-listen-addr=:3000"
  - "-log-format=json"
  - "-log-level=debug"
  - "-otel-enabled"

env:
  - name: "OTEL_SERVICE_NAME"
    value: "my-go-hello-world"
  - name: "OTEL_EXPORTER_OTLP_ENDPOINT"
    value: "http://otel-collector:4318"

annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "3000"
```

Install with custom values:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  -f values-custom.yaml
```

## Application Arguments

The application supports the following command-line flags:

- `-listen-addr` - Server listen address (default `:3000`)
- `-session-key` - Session encryption key
- `-otel-enabled` - Enable OpenTelemetry tracing

### Configuring OpenTelemetry

To enable OpenTelemetry tracing, add the `-otel-enabled` flag and configure OTLP endpoint:

```yaml
args:
  - "-listen-addr=:3000"
  - "-otel-enabled"

env:
  - name: "OTEL_EXPORTER_OTLP_PROTOCOL"
    value: "grpc"
  - name: "OTEL_EXPORTER_OTLP_ENDPOINT"
    value: "http://otel-collector:4317"
  - name: "OTEL_RESOURCE_ATTRIBUTES"
    value: "service.name=go-hello-world,service.version=1.8.2,deployment.environment=production"
```

## Accessing the Application

After installation, you can access the application in several ways:

### Port Forwarding (for testing)

Forward your local port to the service:

```bash
kubectl port-forward svc/go-hello-world 3000:3000
```

Then access the application at `http://localhost:3000` in your browser.

### Via Ingress

Create an Ingress resource to expose the application:

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

Apply with:

```bash
kubectl apply -f ingress.yaml
```

### Via LoadBalancer Service

For cloud environments with LoadBalancer support:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --set service.type=LoadBalancer
```

Then find the external IP:

```bash
kubectl get svc go-hello-world -o wide
```

Access at `http://<EXTERNAL-IP>:3000`

### Via NodePort Service

For on-premise Kubernetes:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --set service.type=NodePort
```

Find the assigned port:

```bash
kubectl get svc go-hello-world -o wide
```

Access at `http://<NODE-IP>:<NODE-PORT>`

## Upgrading

Upgrade to a newer chart version:

```bash
helm upgrade go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
```

Upgrade with custom values:

```bash
helm upgrade go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  -f values-custom.yaml
```

## Uninstalling

Remove the Helm release:

```bash
helm uninstall go-hello-world
```

This will remove all associated Kubernetes resources (deployment, service, etc.).

## Troubleshooting

### Application not accessible

1. Check pod status:
   ```bash
   kubectl get pods -l app.kubernetes.io/name=go-hello-world
   kubectl describe pod <pod-name>
   ```

2. Check logs:
   ```bash
   kubectl logs -l app.kubernetes.io/name=go-hello-world
   ```

3. Verify service is running:
   ```bash
   kubectl get svc go-hello-world
   ```

### Image pull failures

Ensure the image tag is correct and accessible:

```bash
helm show values oci://ghcr.io/wasilak/go-hello-world-chart | grep -A 3 "^image:"
```

### Port already in use

If port 3000 is already in use, specify a different port:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart \
  --set args[0]="-listen-addr=:8080" \
  --set service.targetPort=8080
```

## Prometheus Metrics

The application exposes Prometheus metrics on `/metrics` endpoint. Configure Prometheus scraping:

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "3000"
  prometheus.io/path: "/metrics"
```

## Health Checks

The application responds to HTTP requests at the root path (`/`). Configure Kubernetes health checks:

```yaml
livenessProbe:
  httpGet:
    path: /
    port: 3000
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /
    port: 3000
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Support

- 📖 [Main Project README](https://github.com/wasilak/go-hello-world)
- 🐛 [Issue Tracker](https://github.com/wasilak/go-hello-world/issues)
- 🏢 [GitHub Container Registry](https://ghcr.io/wasilak/go-hello-world)

## Maintainers

| Name | Email |
|------|-------|
| Piotr Boruc | piotr.m.boruc@gmail.com |

---

Chart generated and maintained for go-hello-world reference implementation.
