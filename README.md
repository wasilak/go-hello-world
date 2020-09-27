```
GOOS=linux GOARCH=amd64 go build ./...
```

docker build (multiarch):
```
docker buildx build --tag wasilak/go-hello-world --tag quay.io/wasilak/go-hello-world --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 . --push
```

