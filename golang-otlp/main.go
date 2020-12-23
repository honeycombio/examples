package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	apiTrace "go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/credentials"
)

func main() {
	apikey := flag.String("apikey", "", "Your Honeycomb API Key")
	dataset := flag.String("dataset", "golang-otlp", "Your Honeycomb dataset")
	flag.Parse()

	if *apikey == "" {
		log.Panicln("Honeycomb API key is required. Provide using '--apikey <key>' command line argument.")
	}

	log.Println("Sending trace to dataset: " + *dataset)

	ctx := context.Background()
	exporter, _ := otlp.NewExporter(
		ctx,
		otlp.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		otlp.WithAddress("api.honeycomb.io:443"),
		otlp.WithHeaders(map[string]string{
			"x-honeycomb-team":    *apikey,
			"x-honeycomb-dataset": *dataset,
		}),
	)
	otel.SetTracerProvider(
		sdkTrace.NewTracerProvider(
			sdkTrace.WithConfig(sdkTrace.Config{DefaultSampler: sdkTrace.AlwaysSample()}),
			sdkTrace.WithSpanProcessor(sdkTrace.NewBatchSpanProcessor(exporter)),
			sdkTrace.WithResource(resource.NewWithAttributes(
				semconv.ServiceNameKey.String("golang-otlp"),
				semconv.ServiceVersionKey.String("0.1"),
			)),
		),
	)

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := apiTrace.SpanFromContext(ctx)
		span.SetAttributes(label.String("foo", "bar"))

		_, _ = io.WriteString(w, "Hello, world!\n")
	}

	log.Println("listening at http://localhost:8080")
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(helloHandler), ""))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
