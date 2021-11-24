# Honeycomb OpenTelemetry OTLP example

This is a simple example Sinatra web server that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

Environment configuration:
```
export OTEL_EXPORTER_OTLP_ENDPOINT="https://api.honeycomb.io"
export OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=your-api-key,x-honeycomb-dataset=ruby-otlp"
export OTEL_SERVICE_NAME="sinatra-app"
```

Run the ruby app:
```
bundle install
ruby app.rb
```

Finally, open `http://localhost:4567` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
