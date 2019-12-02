FROM golang:1-alpine as builder

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/wasilak/go-hello-world/

COPY ./ .

ENV GOOS=linux
ENV GOARCH=amd64
RUN go get ./...
RUN go build ./...

FROM alpine:3

COPY --from=builder /go/src/github.com/wasilak/go-hello-world/go-hello-world /usr/local/bin/go-hello-world

ENV SESSION_KEY=cmRiN3VuaTg2Zm9pZ29peWdp

CMD ["go-hello-world"]
