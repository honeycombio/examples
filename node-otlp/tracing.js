"use strict";

const grpc = require("@grpc/grpc-js");
const { Resource } = require("@opentelemetry/resources");
const { ResourceAttributes } = require("@opentelemetry/semantic-conventions");
const { ExpressInstrumentation } = require("@opentelemetry/instrumentation-express");
const { HttpInstrumentation } = require("@opentelemetry/instrumentation-http");
const { NodeTracerProvider } = require("@opentelemetry/node");
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { SimpleSpanProcessor } = require("@opentelemetry/tracing");
const { CollectorTraceExporter } = require("@opentelemetry/exporter-collector-grpc");

const provider = new NodeTracerProvider({
  resource: new Resource({
    [ResourceAttributes.SERVICE_NAME]: 'node-otlp',
  }),
})

const metadata = new grpc.Metadata();
metadata.set('x-honeycomb-team', '<YOUR-APIKEY>');
metadata.set('x-honeycomb-dataset', '<YOUR-DATASET>');

provider.addSpanProcessor(
  new SimpleSpanProcessor(
    new CollectorTraceExporter({
      url: 'grpc://api.honeycomb.io:443/',
      credentials: grpc.credentials.createSsl(),
      metadata
    })
  )
);
provider.register();

registerInstrumentations({
  instrumentations: [
    HttpInstrumentation,
    ExpressInstrumentation
  ]
});

console.log("tracing initialized");
