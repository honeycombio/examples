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
| [`buildevents`](https://github.com/honeycombio/buildevents) | Manual tracing | A small binary to help create traces out of Travis-CI builds, visualizing each step and command as spans within the trace. |
| [`dotnet-core-webapi`](dotnet-core-webapi) | Custom instrumentation with [.NET](https://docs.honeycomb.io/instrumenting-your-application/dotnet/) | A barebones web API with a `HoneycombMiddleware.cs` for capturing HTTP metadata. |
| [`golang-gatekeeper`](golang-gatekeeper) | [Go Beeline](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) and custom instrumentation | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/) |
| [`golang-ratelimiting-proxy`](golang-ratelimiting-proxy) | [Go Beeline](https://docs.honeycomb.io/getting-data-in/beelines/go-beeline/) and cross-service tracing | A rate limiting, tarpitting proxy, intended to be put in front of a web server. |
| [`golang-webapp`](golang-webapp) | [libhoney-go](https://docs.honeycomb.io/sdk/go/) | A two-tier web application (Go+MySQL) which is a Twitter clone. |
| [`golang-wiki-tracing`](golang-wiki-tracing) | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Go | A simple wiki (Go) manually instrumented for tracing. |
| [`honeytail-dockerd`](honeytail-dockerd) | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files) | Using [Honeytail]()'s `keyval` parser to ingest the structured logs of the [Docker]() container engine daemon. |
| [`honeytail-mysql`](honeytail-mysql) | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files) | Using Honeytail's `mysql` parser to ingest MySQL slow query logs |
| [`honeytail-nginx`](honeytail-nginx) | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files) | Using [Honeytail]()'s `nginx` parser to ingest [Nginx]() access logs from an instance acting as a reverse proxy. |
| [`java-beeline`](java-beeline) | [beeline-java](https://docs.honeycomb.io/getting-data-in/java/beeline/)| A simple web app instrumented for tracing with the Java Beeline for SpringBoot |
| [`java-webapp`](java-webapp) | [libhoney-java](https://docs.honeycomb.io/sdk/java/) | A TODO API written and instrumented using Java Spring |
| [`kubernetes-envoy-tracing`](kubernetes-envoy-tracing) | Using the [Honeycomb Opentracing Proxy](https://github.com/honeycombio/honeycomb-opentracing-proxy) to accept OpenTracing data | Two small services deployed to Kubernetes which communicate using [Envoy Proxy](https://www.envoyproxy.io/) |
| [`node-tracing-example`](node-tracing-example) | [Node Beeline](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/) | A simple webapp showing intra-service and cross-service tracing. |
| [`node-serverless-app`](node-serverless-app) | [Node Beeline](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/) + Lambda | A simple Lambda function meant to be part of a larger trace. |
| [`python-api`](python-api) | [libhoney-py](https://docs.honeycomb.io/sdk/python/) | A TODO API written and instrumented using Python (Flask). |
| [`python-gatekeeper`](python-gatekeeper) | [Python Beeline](https://docs.honeycomb.io/getting-data-in/beelines/python-beeline/) and custom instrumentation | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/) |
| [`ruby-gatekeeper`](ruby-gatekeeper) | [Ruby Beeline](https://docs.honeycomb.io/getting-data-in/beelines/ruby-beeline/) and custom instrumentation | An API server paired with the [Gatekeeper Tour](https://docs.honeycomb.io/gatekeeper-tour/) |
| [`ruby-wiki-tracing`](ruby-wiki-tracing) | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Ruby | A simple wiki (Ruby) manually instrumented for tracing. |
| [`webhook-listener-triggers`](webhook-listener-triggers) | Executing a webhook as a result of a [Honeycomb Trigger](https://docs.honeycomb.io/working-with-data/triggers/) firing | A small Go application which listens for HTTP requests issued as a result of a trigger firing |


## Proposed Examples

The following have been proposed but not implemented:

| Directory           | Description                                                             |
| ------------------- | ----------------------------------------------------------------------- |
| `javascript-api`    | A TODO API written and instrumented using JavaScript.                   |
| `sidekiq`           | Observing behavior of the background job runner Sidekiq.                |
| `honeytail-apache`  | Ingesting Apache access logs using Honeytail.                           |
| `honeytail-haproxy` | Ingesting HAProxy access logs using Honeytail.                          |
| `logstash`          | Using the Honeycomb Logstash plugin to send parsed events to Honeycomb. |
| `fluentd`           | Using the Honeycomb Fluentd plugin to send parsed events to Honeycomb.  |

We highly encourage community contribution! Let us know if there's anything you'd like to see
by [filing an issue](https://github.com/honeycombio/examples/issues/new) and CC-ing Honeycombers
for discussion.

Let us know if there is something specific you'd like to see by [filing an
issue](https://github.com/honeycombio/examples/issues/new).
