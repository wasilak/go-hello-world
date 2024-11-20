FROM quay.io/wasilak/golang:1.23 AS builder

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o /go-hello-world

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/wasilak/go-hello-world"

COPY --from=builder /go-hello-world .

ENV USER=root

ENTRYPOINT ["/go-hello-world"]
