FROM golang:alpine

RUN apk add --update --no-cache git
RUN go get github.com/honeycombio/honeytail

FROM alpine

RUN apk add --update --no-cache ca-certificates
COPY --from=0 /go/bin/honeytail /usr/bin/honeytail
ENTRYPOINT ["/usr/bin/honeytail"]
