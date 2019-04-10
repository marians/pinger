FROM golang:1.12-alpine as builder
COPY . /go/src/github/marians/pinger
WORKDIR /go/src/github/marians/pinger
RUN go build

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github/marians/pinger/pinger /pinger
ENTRYPOINT ["/pinger"]
