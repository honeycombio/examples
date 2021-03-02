
# Honeycomb OpenTelemetry OTLP example

This is a simple Python flask web app that uses OpenTelemetry to generate trace data and send to Honeycomb using the OTLP exporter.

First you need to edit [app.py](./app.app.py) to set your Honeycomb API key and dataset.

Then run the web app.
```
python app.py
```

Finally, open `http://localhost:5000` to generate some trace data that will be visble in the [Honeycomb UI](http://ui.honeycomb.io).
