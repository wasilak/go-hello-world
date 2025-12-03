![icon_small](icon_small.jpg)

# go-hello-world

[![Maintainability](https://api.codeclimate.com/v1/badges/7ada7f029e74c805ec1c/maintainability)](https://codeclimate.com/github/wasilak/go-hello-world/maintainability)

Simple web app for testing e.g. container orchestrators, autoscalers etc.

```
GOOS=linux GOARCH=amd64 go build ./...
```

docker build (multiarch):
```
# setup, if needed
docker buildx create --use --append --name mybuilder unix:///var/run/docker.sock

# actual build
docker buildx build --tag quay.io/wasilak/go-hello-world --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 . --push
```

## Helm Installation

Deploy to Kubernetes using the Helm chart from OCI registry:

```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
```

For detailed installation instructions, configuration options, and usage examples, see the [Helm Chart Documentation](./charts/go-hello-world/README.md).

### Quick Access Methods

After installation:

```bash
# Port-forward to access locally
kubectl port-forward svc/go-hello-world 3000:3000
```

Access at `http://localhost:3000`

For more information on Kubernetes deployment, ingress configuration, and troubleshooting, refer to the [chart README](./charts/go-hello-world/README.md).

### Migration from GitHub Pages Repository

The go-hello-world Helm chart has transitioned from GitHub Pages (chart-releaser) distribution to OCI registry (GHCR). This provides better integration with modern Helm tooling and container registry infrastructure.

**Phase 1 (v1.9.0):** Both distribution methods are active. You can use either approach:

**Old method (GitHub Pages) - Deprecated:**
```bash
# Add the Helm repository (no longer recommended)
helm repo add go-hello-world https://wasilak.github.io/go-hello-world
helm repo update
helm install go-hello-world go-hello-world/go-hello-world
```

**New method (OCI Registry) - Recommended:**
```bash
helm install go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
```

**Phase 2 (v2.0.0+):** Only OCI registry distribution. GitHub Pages repository is deprecated.

**Migration Steps:**

If you're currently using the GitHub Pages method:

1. Remove the old repository:
   ```bash
   helm repo remove go-hello-world
   ```

2. Upgrade using the OCI registry:
   ```bash
   helm upgrade go-hello-world oci://ghcr.io/wasilak/go-hello-world-chart
   ```

**Why OCI Registry?**
- Direct installation without repository management
- Better integration with container registries
- Simpler discovery and management
- Consistent with Helm 3+ best practices

How to run (example):

```bash
OTEL_EXPORTER_OTLP_PROTOCOL=grpc OTEL_RESOURCE_ATTRIBUTES="service.name=go-hello-world,service.version=v0.0.6,deployment.environment=test" OTEL_EXPORTER_OTLP_ENDPOINT="https://localhost:4317" OTEL_METRICS_EXPORTER="otlp" OTEL_LOGS_EXPORTER="otlp" go run . -listen-addr :3000 -session-key f98ehv9273hvreof -otel-enabled
```
