require 'sinatra/base'
require 'opentelemetry/sdk'
require 'opentelemetry/exporter/otlp'


# configure SDK with OTLP exporter
OpenTelemetry::SDK.configure do |c|
  c.add_span_processor(
    OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
      exporter: OpenTelemetry::Exporter::OTLP::Exporter.new(
		endpoint: 'localhost:55681/v1/trace', # send to local collector
		insecure: true
      )
    )
  )
end

class OpenTelemetryMiddleware
	def initialize(app)
	  @app = app
	  @tracer = OpenTelemetry.tracer_provider.tracer('sinatra', '1.0')
	end

	def call(env)
		# Extract context from request headers
		context = OpenTelemetry.propagation.http.extract(env)

		status, headers, response_body = 200, {}, ''

		# Span name SHOULD be set to route:
		span_name = env['PATH_INFO']

		# For attribute naming, see
		# https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-semantic-conventions.md#http-server

		# Span kind MUST be `:server` for a HTTP server span
		@tracer.in_span(
		  span_name,
		  attributes: {
			'component' => 'http',
			'http.method' => env['REQUEST_METHOD'],
			'http.route' => env['PATH_INFO'],
			'http.url' => env['REQUEST_URI'],
		  },
		  kind: :server,
		  with_parent: context
		) do |span|
		  # Run application stack
		  status, headers, response_body = @app.call(env)

		  span.set_attribute('http.status_code', status)
		end

		[status, headers, response_body]
	  end
end

class App < Sinatra::Base
  set :bind, '0.0.0.0'
  use OpenTelemetryMiddleware

  get '/' do
	'Hello world!'
  end

  run! if app_file == $0
end
