FROM golang:alpine

ADD . /go/src/github.com/honeycombio/examples/golang
RUN go install github.com/honeycombio/examples/golang

FROM alpine

RUN apk add --update --no-cache ca-certificates
COPY --from=0 /go/bin/golang /usr/bin/shoutr
ENTRYPOINT ["/usr/bin/shoutr"]
