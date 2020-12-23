# Honeycomb OpenTelemetry OTLP example

This is a simple NodeJS express web app that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

First you need to edit [tracing.js](./tracing.js) to set your Honeycomb API key and dataset.

Then run the web app.
```
npm install
npm start
```

Finally, open `http://localhost:8080` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
