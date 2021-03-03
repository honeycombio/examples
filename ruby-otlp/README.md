# Honeycomb OpenTelemetry OTLP example

This is a simple example Sinatra web server that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

The [Ruby OpenTelemetry OTLP exporter](https://github.com/open-telemetry/opentelemetry-ruby/tree/master/exporter/otlp) currently only works over HTTP with JSON and the Honeycomb ingest API only works over gRPC. This means we need to setup and run an [OpenTelemetry Collector]() to receive the trace data over JSON and then send to Honeycomb over gRPC.

The simplest way to run a collector is to use the offical OpenTelemetry Collector Docker image and pass in a configuration file. This [example config](./collector-config.yaml) can be used to setup a minimal collector that receives OTLP data over both HTTP and gRPC and forwards onto Honeycomb using gRPC. You will need to edit the config to provide your Honeycomb API Key and dataset name. 

Then run the following docker command to start the collector:
```bash
docker run --rm -p 4317:4317 -p 55680-55681:55680-55681 -v "${PWD}/collector-config.yaml":/collector-config.yaml otel/opentelemetry-collector --config collector-config.yaml
```

Export an environment variable for the OTLP exporter to send data to the collector:
```
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:55681
```

Run the ruby app:
```
bundle install
ruby app.rb
```

Finally, open `http://localhost:4567` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
