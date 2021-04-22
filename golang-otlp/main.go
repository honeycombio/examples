package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	apiTrace "go.opentelemetry.io/otel/trace"
)

func initTracer() func() {
	// Fetch the necessary settings (from environment variables, in this example).
	// You can find the API key via https://ui.honeycomb.io/account after signing up for Honeycomb.
	apikey, _ := os.LookupEnv("HONEYCOMB_API_KEY")
	dataset, _ := os.LookupEnv("HONEYCOMB_DATASET")
	if apikey == "" {
		log.Panicln("Honeycomb API key is required. Set the HONEYCOMB_API_KEY environment variable.")
	}
	if dataset == "" {
		dataset = "golang-otlp"
	}
	log.Println("Sending trace to dataset: " + dataset)

	// Initialize an OTLP exporter over gRPC and point it to Honeycomb.
	ctx := context.Background()
	exporter, err := otlp.NewExporter(
		ctx,
		otlpgrpc.NewDriver(
			// otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
			otlpgrpc.WithInsecure(),
			otlpgrpc.WithEndpoint("localhost:4317"),
			// otlpgrpc.WithHeaders(map[string]string{
			// 	"x-honeycomb-team":    apikey,
			// 	"x-honeycomb-dataset": dataset,
			// }),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the OTel tracer provider.
	provider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(&BaggageSpanProcessor{}),
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String("golang-otlp"),
			semconv.ServiceVersionKey.String("0.1"),
		)),
	)
	otel.SetTracerProvider(provider)

	// This callback will ensure all spans get flushed before the program exits.
	return func() {
		ctx := context.Background()
		err := provider.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	cleanup := initTracer()
	defer cleanup()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := apiTrace.SpanFromContext(ctx)
		defer span.End()
		span.SetAttributes(attribute.String("foo", "bar"))
		addTraceField(ctx, span, attribute.String("username", "mike"))

		testSpan := apiTrace.SpanFromContext(ctx)
		defer testSpan.End()
		testSpan.SetAttributes(attribute.String("bar", "baz"))

		_, _ = io.WriteString(w, "Hello, world!\n")
	}

	log.Println("listening at http://localhost:8080")
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(helloHandler), ""))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func addTraceField(ctx context.Context, span apiTrace.Span, attr attribute.KeyValue) {
	span.SetAttributes(attr)
	ctx = baggage.ContextWithValues(ctx, attr)
}

type BaggageSpanProcessor struct {
}

func (ssp *BaggageSpanProcessor) OnStart(ctx context.Context, span sdkTrace.ReadWriteSpan) {
	baggageVals := baggage.Set(ctx)
	span.SetAttributes(baggageVals.ToSlice()...)
}

func (ssp *BaggageSpanProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
}

func (ssp *BaggageSpanProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (ssp *BaggageSpanProcessor) ForceFlush(context.Context) error {
	return nil
}
