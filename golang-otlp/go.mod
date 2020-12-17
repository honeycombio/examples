module github.com/honeycombio/examples/golang-otlp

go 1.14

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.1
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/otlp v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	google.golang.org/grpc v1.32.0
)
