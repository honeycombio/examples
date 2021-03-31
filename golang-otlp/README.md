
# Honeycomb OpenTelemetry OTLP example

This is a simple golang web app that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

Run the web app and pass in your API key and dataset as environment variables:
```
HONEYCOMB_API_KEY=<your-api-key> HONEYCOMB_DATASET=<your-dataset> go run main.go
```

Next, open `http://localhost:8080` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
