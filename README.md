# Honeycomb Instrumentation Examples

The full power of Honeycomb is unlocked by proper instrumentation, especially at
the code level. The examples in this repository will help guide you
along the way to unlocking this power for yourself by showing you how
instrumentation is meant to be done.

Most of the top level directories in this repository correspond to an example
which demonstrates instrumentation using Honeycomb and sample queries to help
you along. The current examples are:

| Directory | Meant to Teach | Description |
| --- | --- | --- |
| `golang-gatekeeper` | [Go Beeline](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) | An API server paired with the [gatekeeper tour](https://ui.honeycomb.io/quickstart/datasets/gatekeeper-tour) |
| `ruby-gatekeeper` | [Ruby Beeline](https://docs.honeycomb.io/getting-data-in/beelines/ruby-beeline/) | An API server paired with the [gatekeeper tour](https://ui.honeycomb.io/quickstart/datasets/gatekeeper-tour) |
| `python-gatekeeper` | [Python Beeline](https://docs.honeycomb.io/getting-data-in/beelines/python-beeline/) | An API server paired with the [gatekeeper tour](https://ui.honeycomb.io/quickstart/datasets/gatekeeper-tour) |
| `golang-webapp` | [libhoney-go](https://docs.honeycomb.io/sdk/go/) | A two-tier web application (Go+MySQL) which is a Twitter clone. |
| `honeytail-dockerd` | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files) | Using [Honeytail]()'s `keyval` parser to ingest the structured logs of the [Docker]() container engine daemon. |
| `honeytail-nginx` | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files) | Using [Honeytail]()'s `nginx` parser to ingest [Nginx]() access logs from an instance acting as a reverse proxy. |
| `python-api` | [libhoney-py](https://docs.honeycomb.io/sdk/python/) | A TODO API written and instrumented using Python (Flask). |
| `golang-wiki-tracing` | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Go | A simple wiki (Go) manually instrumented for tracing. |
| `ruby-wiki-tracing` | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Ruby | A simple wiki (Ruby) manually instrumented for tracing. |
| `java-webapp` | [libhoney-java](https://docs.honeycomb.io/sdk/java/) | A TODO API written and instrumented using Java Spring |
| `kubernetes-envoy-tracing` | Using the [Honeycomb Opentracing Proxy](https://github.com/honeycombio/honeycomb-opentracing-proxy) to accept OpenTracing data | Two small services deployed to Kubernetes which communicate using [Envoy Proxy](https://www.envoyproxy.io/) |
| `webhook-listener-triggers` | Executing a webhook as a result of a [Honeycomb Trigger](https://docs.honeycomb.io/working-with-data/triggers/) firing | A small Go application which listens for HTTP requests issued as a result of a trigger firing |

## Proposed Examples

The following have been proposed but not implemented:

| Directory | Description |
| --- | --- |
| `javascript-api` | A TODO API written and instrumented using JavaScript. |
| `sidekiq` | Observing behavior of the background job runner Sidekiq. |
| `honeytail-apache` | Ingesting Apache access logs using Honeytail. |
| `honeytail-haproxy` | Ingesting HAProxy access logs using Honeytail. |
| `logstash` | Using the Honeycomb Logstash plugin to send parsed events to Honeycomb. |
| `fluentd` | Using the Honeycomb Fluentd plugin to send parsed events to Honeycomb. |

We highly encourage community contribution! Let us know if there's anything you'd like to see
by [filing an issue](https://github.com/honeycombio/examples/issues/new) and CC-ing Honeycombers
for discussion.

Let us know if there is something specific you'd like to see by [filing an
issue](https://github.com/honeycombio/examples/issues/new).
