FROM golang:alpine

ADD . /go/src/github.com/honeycombio/examples/golang
RUN go install github.com/honeycombio/examples/golang

FROM alpine

RUN apk add --update --no-cache ca-certificates
RUN mkdir -p /opt/shoutr/bin
WORKDIR /opt/shoutr
COPY --from=0 /go/bin/golang /opt/shoutr/shoutr
COPY templates /opt/shoutr/templates
ENTRYPOINT ["/opt/shoutr/shoutr"]
