module github.com/honeycombio/examples/golang-otlp

go 1.14

require (
	github.com/golang/protobuf v1.4.3 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/otlp v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
	google.golang.org/grpc v1.36.0
)
