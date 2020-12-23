
# Honeycomb OpenTelemetry OTLP example

This is a simple golang web app that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

Run the web app and pass in youe API key and dataset:
```
go run main.go -apikey <your-apikey> -dataset <dataset-name>
```

Next, open `http://localhost:8080` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
