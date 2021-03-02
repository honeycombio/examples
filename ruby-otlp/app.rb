require 'opentelemetry/sdk'
require 'opentelemetry/exporter/otlp'
require 'rubygems'
require 'bundler/setup'

Bundler.require

OpenTelemetry::SDK.configure do |c|
  c.service_name = 'ruby-otlp'
  c.use 'OpenTelemetry::Instrumentation::Sinatra'
  c.add_span_processor(
    OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
      OpenTelemetry::Exporter::OTLP::Exporter.new(
        endpoint: 'http://localhost:55681/v1/traces'
      )
    )
  )
end

get '/' do
  'Hello world!'
end
