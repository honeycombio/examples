"use strict";

const grpc = require("@grpc/grpc-js");
const { ExpressInstrumentation } = require("@opentelemetry/instrumentation-express");
const { HttpInstrumentation } = require("@opentelemetry/instrumentation-http");
const { NodeTracerProvider } = require("@opentelemetry/node");
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { SimpleSpanProcessor } = require("@opentelemetry/tracing");
const { CollectorTraceExporter } = require("@opentelemetry/exporter-collector-grpc");

const metadata = new grpc.Metadata();
metadata.set('x-honeycomb-team', '<YOUR-APIKEY>');
metadata.set('x-honeycomb-dataset', '<YOUR-DATASET>');

const provider = new NodeTracerProvider();
provider.addSpanProcessor(
  new SimpleSpanProcessor(
    new CollectorTraceExporter({
      serviceName: 'node-otlp',
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
