# Honeycomb OpenTelemetry OTLP example

This is a simple example Sinatra web server that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

The [Ruby OpenTelemetry OTLP exporter](https://github.com/open-telemetry/opentelemetry-ruby/tree/master/exporter/otlp) currently only works over HTTP with JSON and the Honeycomb ingest API only works over gRPC. This means we need to
setup and run a collector to receive the trace data over JSON and then send to Honeycomb over gRPC.

The simplest way to run a collector is to use the offical OpenTelemetry Collector Docker image and pass in the configuration file. First edit the config to provide your Honeycomb API Key and dataset name. Then run the following to start the collector:

```
docker run --rm -p 55680-55681:55680-55681 -v "${PWD}/collector-config.yaml":/collector-config.yaml otel/opentelemetry-collector --config collector-config.yaml
```

Once the collector is up and running, you then can run the ruby app with the following:
```
bundle install
ruby app.rb
```

Finally, open `http://localhost:4567` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
