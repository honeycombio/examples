"use strict";

const grpc = require('grpc');
const { LogLevel } = require("@opentelemetry/core");
const { NodeTracerProvider } = require("@opentelemetry/node");
const { SimpleSpanProcessor } = require("@opentelemetry/tracing");
const { CollectorTraceExporter } = require("@opentelemetry/exporter-collector-grpc");

const provider = new NodeTracerProvider({ logLevel: LogLevel.ERROR });

const metadata = new grpc.Metadata();
metadata.set("x-honeycomb-team", "");
metadata.set("x-honeycomb-dataset", "");

const collectorOptions = {
	serviceName: 'node-otlp',
	url: 'api.honeycomb.io:443',
	credentials: grpc.credentials.createSsl(),
	metadata
  };

provider.addSpanProcessor(
  new SimpleSpanProcessor(
    new CollectorTraceExporter(collectorOptions)
  )
);

provider.register();
console.log("tracing initialized");
