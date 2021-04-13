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

