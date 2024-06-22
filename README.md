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

How to run (example):

```bash
OTEL_EXPORTER_OTLP_PROTOCOL=grpc OTEL_RESOURCE_ATTRIBUTES="service.name=go-hello-world,service.version=v0.0.6,deployment.environment=test" OTEL_EXPORTER_OTLP_ENDPOINT="https://localhost:4317" OTEL_METRICS_EXPORTER="otlp" OTEL_LOGS_EXPORTER="otlp" go run . -listen-addr :3000 -session-key f98ehv9273hvreof -otel-enabled
```
