# golang-ratelimiting-proxy

This example contains a simple HTTP proxy which rate limits based on remote IP address before passing the requests on to a localhost webserver.

This is a **fully-instrumented** service using the [Beeline for Go](https://github.com/honeycombio/beeline-go/) that propagates tracing headers through to the downstream application.

## More examples

See our announcement of the [Beeline for Go v2](https://www.honeycomb.io/blog/2018/09/the-honeycomb-beeline-for-go-v2-is-go/) for more information and discussion of how to trace across service boundaries.
