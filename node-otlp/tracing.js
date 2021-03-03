"use strict";

const grpc = require('grpc');
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
			url: 'api.honeycomb.io:443',
			credentials: grpc.credentials.createSsl(),
			metadata
		})
  )
);
provider.register();

registerInstrumentations({
  tracerProvider: provider,
});

console.log("tracing initialized");
