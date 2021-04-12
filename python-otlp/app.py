import flask
import requests

from grpc import ssl_channel_credentials

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.flask import FlaskInstrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

otlp_exporter = OTLPSpanExporter(
	endpoint="api.honeycomb.io:443",
	credentials=ssl_channel_credentials(),
	headers=(("x-honeycomb-team", ""),("x-honeycomb-dataset","python-otlp"))
)

trace.set_tracer_provider(TracerProvider(resource=Resource({"service.name": "python-otlp", "service.version":"0.1"})))
trace.get_tracer_provider().add_span_processor(
    BatchSpanProcessor(otlp_exporter)
)

app = flask.Flask(__name__)
FlaskInstrumentor().instrument_app(app)
RequestsInstrumentor().instrument()


@app.route("/")
def hello():
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("example-request"):
        requests.get("http://www.example.com")
    return "hello"


app.run(debug=True, port=5000)
