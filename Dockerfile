FROM  --platform=$BUILDPLATFORM quay.io/wasilak/golang:1.20-alpine as builder
COPY --from=tonistiigi/xx:golang / /

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/wasilak/go-hello-world/

COPY ./ .

RUN go build .

FROM --platform=$BUILDPLATFORM quay.io/wasilak/alpine:3

COPY --from=builder /go/src/github.com/wasilak/go-hello-world/go-hello-world /usr/local/bin/go-hello-world

ENV SESSION_KEY=cmRiN3VuaTg2Zm9pZ29peWdp

CMD ["/usr/local/bin/go-hello-world"]
