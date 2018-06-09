from flask import Flask, request, abort, jsonify
import socket
import os
import sys
import requests
import json
import time
import random

app = Flask(__name__)

TRACE_HEADERS_TO_PROPAGATE = [
    "X-Ot-Span-Context",
    "X-Request-Id",

    # Zipkin headers
    "X-B3-TraceId",
    "X-B3-SpanId",
    "X-B3-ParentSpanId",
    "X-B3-Sampled",
    "X-B3-Flags",

    # Jaeger header (for native client)
    "uber-trace-id"
]

@app.route("/stage/<int:i>")
def index(i):
    if i == 3:
        time.sleep(0.3 + (3.0*random.random()))
    return jsonify({"upstream": "healthy"})

@app.route("/echo/<name>")
def trace(name):
    headers = {}
    if name == "chris":
        print(json.dumps({"error": "too many chrises"}))
        abort(500)

    # call service 2 from service 1
    if int(os.environ["SERVICE_NAME"]) == 1 :
        for header in TRACE_HEADERS_TO_PROPAGATE:
            if header in request.headers:
                headers[header] = request.headers[header]
        for i in range(0, 5):
            ret = requests.get("http://localhost:9000/stage/{}".format(i), headers=headers)
            print(json.dumps({"upstream_response": str(ret), name: "name"}))
    else:
        return jsonify({"upstream": "healthy"})
    print(json.dumps({"normal": True}))
    return ("""<img src="https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy.png" height="200" />
<pre><code>Hello {}!

Served by:
service {}
pod: {}
resolved hostname: {}</code></pre>
""".format(name,
           os.environ["SERVICE_NAME"],
           socket.gethostname(),
           socket.gethostbyname(socket.gethostname())))

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080, debug=True)
