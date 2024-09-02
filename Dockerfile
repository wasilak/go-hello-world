FROM  quay.io/wasilak/golang:1.23-alpine as builder

LABEL org.opencontainers.image.source="https://github.com/wasilak/go-hello-world"

RUN apk add --no-cache git

WORKDIR /src

COPY ./ .

RUN go build .

FROM quay.io/wasilak/alpine:3

COPY --from=builder /src/go-hello-world /bin/go-hello-world

ENV SESSION_KEY=cmRiN3VuaTg2Zm9pZ29peWdp

CMD ["/bin/go-hello-world"]
