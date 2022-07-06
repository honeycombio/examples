# Honeycomb Instrumentation Examples

[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/examples?color=success)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)

The full power of Honeycomb is unlocked by proper instrumentation, especially at
the code level. The examples in this repository will help guide you
along the way to unlocking this power for yourself by showing you how
instrumentation is meant to be done.

Most of the top level directories in this repository correspond to an example
which demonstrates instrumentation using Honeycomb and sample queries to help
you along. The current examples are:

| Examples                                                                                                             | Meant to Teach                                                                                                                 | Description                                                                                                                                |
|----------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| [`buildevents`](https://github.com/honeycombio/buildevents)                                                          | Manual tracing                                                                                                                 | A small binary to help create traces out of Travis-CI builds, visualizing each step and command as spans within the trace.                 |
| [`golang-gatekeeper`](golang-gatekeeper)                                                                             | [Go Beeline](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) and custom instrumentation                        | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/)                                                |
| [`golang-ratelimiting-proxy`](golang-ratelimiting-proxy)                                                             | [Go Beeline](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) and cross-service tracing                         | A rate limiting, tarpitting proxy, intended to be put in front of a web server.                                                            |
| [`golang-webapp`](golang-webapp)                                                                                     | [libhoney-go](https://docs.honeycomb.io/sdk/go/)                                                                               | A two-tier web application (Go+MySQL) which is a Twitter clone.                                                                            |
| [`golang-wiki-tracing`](golang-wiki-tracing)                                                                         | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Go                  | A simple wiki (Go) manually instrumented for tracing.                                                                                      |
| [`honeytail`](https://github.com/honeycombio/honeytail/tree/main/example)                                            | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files)                                             | Using [Honeytail]()'s `keyval` parser to ingest a structured log file.                                                                     |
| [`honeytail-haproxy`](honeytail-haproxy)                                                                             | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files)                                             | Using [Honeytail]()'s `nginx` parser to ingest [HAProxy](https://www.haproxy.org/) access logs from an instance acting as a reverse proxy. |
| [`honeytail-mysql`](honeytail-mysql)                                                                                 | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files)                                             | Using Honeytail's `mysql` parser to ingest MySQL slow query logs                                                                           |
| [`honeytail-nginx`](honeytail-nginx)                                                                                 | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files)                                             | Using [Honeytail]()'s `nginx` parser to ingest [Nginx]() access logs from an instance acting as a reverse proxy.                           |
| [`honeycomb-opentelemetry-java`](https://github.com/honeycombio/honeycomb-opentelemetry-java/tree/main/examples)     | [OpenTelemetry Java](https://docs.honeycomb.io/getting-data-in/opentelemetry/java-distro/)                                     | Spring Boot applications instrumented using the Honeycomb OpenTelemetry Java Distribution                                                  |
| [`honeycomb-opentelemetry-dotnet`](https://github.com/honeycombio/honeycomb-opentelemetry-dotnet/tree/main/examples) | [OpenTelemetry .NET](https://docs.honeycomb.io/getting-data-in/opentelemetry/dotnet-distro/)                                   | .NET applications instrumented using the Honeycomb OpenTelemetry .NET Distribution                                                         | 
| [`libhoney-java`](https://github.com/honeycombio/libhoney-java/tree/main/examples)                                   | [Libhoney Java](https://docs.honeycomb.io/getting-data-in/libhoney/java/)                                                      | Examples of using the low-level Libhoney Java SDK                                                                                          |
| [`kubernetes-envoy-tracing`](kubernetes-envoy-tracing)                                                               | Using the [Honeycomb Opentracing Proxy](https://github.com/honeycombio/honeycomb-opentracing-proxy) to accept OpenTracing data | Two small services deployed to Kubernetes which communicate using [Envoy Proxy](https://www.envoyproxy.io/)                                |
| [`beeline-nodejs`](https://github.com/honeycombio/beeline-nodejs/tree/main/examples/node-tracing)                    | [Node Beeline](https://docs.honeycomb.io/getting-data-in/beeline/nodejs/)                                                      | A simple webapp instrumented with the NodeJS Beeline.                                                                                      |
| [`node-serverless-app`](node-serverless-app)                                                                         | [Node Beeline](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/) + Lambda                                  | A simple Lambda function meant to be part of a larger trace.                                                                               |
| [`python-api`](python-api)                                                                                           | [libhoney-py](https://docs.honeycomb.io/sdk/python/)                                                                           | A TODO API written and instrumented using Python (Flask).                                                                                  |
| [`python-gatekeeper`](python-gatekeeper)                                                                             | [Python Beeline](https://docs.honeycomb.io/getting-data-in/beelines/python-beeline/) and custom instrumentation                | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/)                                                |
| [`ruby-gatekeeper`](ruby-gatekeeper)                                                                                 | [Ruby Beeline](https://docs.honeycomb.io/getting-data-in/beelines/ruby-beeline/) and custom instrumentation                    | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/)                                                |
| [`ruby-wiki-tracing`](ruby-wiki-tracing)                                                                             | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Ruby                | A simple wiki (Ruby) manually instrumented for tracing.                                                                                    |
| [`webhook-listener-triggers`](webhook-listener-triggers)                                                             | Executing a webhook as a result of a [Honeycomb Trigger](https://docs.honeycomb.io/working-with-data/triggers/) firing         | A small Go application which listens for HTTP requests issued as a result of a trigger firing                                              |


## Proposed Examples

The following have been proposed but not implemented:

| Directory           | Description                                                             |
| ------------------- | ----------------------------------------------------------------------- |
| `javascript-api`    | A TODO API written and instrumented using JavaScript.                   |
| `sidekiq`           | Observing behavior of the background job runner Sidekiq.                |
| `honeytail-apache`  | Ingesting Apache access logs using Honeytail.                           |
| `logstash`          | Using the Honeycomb Logstash plugin to send parsed events to Honeycomb. |
| `fluentd`           | Using the Honeycomb Fluentd plugin to send parsed events to Honeycomb.  |

We highly encourage community contribution! Let us know if there's anything you'd like to see
by [filing an issue](https://github.com/honeycombio/examples/issues/new) and CC-ing Honeycombers
for discussion.

Let us know if there is something specific you'd like to see by [filing an
issue](https://github.com/honeycombio/examples/issues/new).
