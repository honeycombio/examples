# java-otel-manual

This is a simple Java web app that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

First you need to edit [JavaOtlpApplication.java](./src/main/java/io/honeycomb/examples/javaotlp/JavaOtlpApplication.java) to set your Honeycomb API key and dataset.

Then run the web app:
```sh
HONEYCOMB_API_ENDPOINT=https://api-dogfood.honeycomb.io:443 \
HONEYCOMB_API_KEY=abc \
HONEYCOMB_DATASET=shiny-java \
./gradlew bootRun
```

Finally, open `http://localhost:8080` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
