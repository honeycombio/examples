FROM golang:alpine

ADD . /go/src/github.com/honeycombio/examples/golang
RUN go install github.com/honeycombio/examples/golang

FROM alpine

COPY --from=0 /go/bin/golang /usr/bin/shoutr
RUN apk add --update --no-cache ca-certificates
ENTRYPOINT ["/usr/bin/shoutr"]
