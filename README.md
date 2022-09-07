# Honeycomb Instrumentation Examples

[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/examples?color=success)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)

The full power of Honeycomb is unlocked by proper instrumentation, especially at
the code level. The examples in this repository will help guide you
along the way to unlocking this power for yourself by showing you how
instrumentation is meant to be done.

Most of the top level directories in this repository correspond to an example
which demonstrates instrumentation using Honeycomb and sample queries to help
you along. The current examples are:

| Examples                                                                                                             | Meant to Teach                                                                                                                 | Description                                                                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`Example Greeting Service`](https://github.com/honeycombio/example-greeting-service)                                | Instrumentation, tracing across services                                                                                       | A ridiculously over-engineered microservice application implemented in several languages. The services are instrumented to send telemetry to Honeycomb, with examples of Beelines and OpenTelemetry SDKs. |
| [`buildevents`](https://github.com/honeycombio/buildevents)                                                          | Manual tracing                                                                                                                 | A small binary to help create traces out of Travis-CI builds, visualizing each step and command as spans within the trace.                                                                                |
| [`libhoney-go`](https://github.com/honeycombio/libhoney-go/tree/main/examples)                                       | [libhoney-go](https://docs.honeycomb.io/getting-data-in/libhoney/go/)                                                          | Examples of using the low-level Libhoney Go SDK, including manual tracing.                                                                                                                                |
| [`honeytail`](https://github.com/honeycombio/honeytail/tree/main/examples)                                           | [Honeytail](https://docs.honeycomb.io/getting-data-in/logs/honeytail/) (flat log files)                                        | Examples of using Honeytail to ingest structured log files.                                                                                                                                               |
| [`honeytail-mysql`](honeytail-mysql)                                                                                 | [Honeytail](https://docs.honeycomb.io/getting-data-in/honeytail/) (flat log files)                                             | Using Honeytail's `mysql` parser to ingest MySQL slow query logs                                                                                                                                          |
| [`honeycomb-opentelemetry-java`](https://github.com/honeycombio/honeycomb-opentelemetry-java/tree/main/examples)     | [OpenTelemetry Java](https://docs.honeycomb.io/getting-data-in/opentelemetry/java-distro/)                                     | Spring Boot applications instrumented using the Honeycomb OpenTelemetry Java Distribution                                                                                                                 |
| [`honeycomb-opentelemetry-dotnet`](https://github.com/honeycombio/honeycomb-opentelemetry-dotnet/tree/main/examples) | [OpenTelemetry .NET](https://docs.honeycomb.io/getting-data-in/opentelemetry/dotnet-distro/)                                   | .NET applications instrumented using the Honeycomb OpenTelemetry .NET Distribution                                                                                                                        | 
| [`honeycomb-opentelemetry-go`](https://github.com/honeycombio/honeycomb-opentelemetry-go/tree/main/examples)         | [OpenTelemetry Go](https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/)                                         | Go applications instrumented using the Honeycomb OpenTelemetry Go Distribution                                                                                                                            | 
| [`libhoney-java`](https://github.com/honeycombio/libhoney-java/tree/main/examples)                                   | [Libhoney Java](https://docs.honeycomb.io/getting-data-in/libhoney/java/)                                                      | Examples of using the low-level Libhoney Java SDK                                                                                                                                                         |
| [`kubernetes-envoy-tracing`](kubernetes-envoy-tracing)                                                               | Using the [Honeycomb Opentracing Proxy](https://github.com/honeycombio/honeycomb-opentracing-proxy) to accept OpenTracing data | Two small services deployed to Kubernetes which communicate using [Envoy Proxy](https://www.envoyproxy.io/)                                                                                               |
| [`beeline-nodejs`](https://github.com/honeycombio/beeline-nodejs/tree/main/examples/node-tracing)                    | [Node Beeline](https://docs.honeycomb.io/getting-data-in/beeline/nodejs/)                                                      | A simple webapp instrumented with the NodeJS Beeline.                                                                                                                                                     |
| [`node-serverless-app`](node-serverless-app)                                                                         | [Node Beeline](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/) + Lambda                                  | A simple Lambda function meant to be part of a larger trace.                                                                                                                                              |
| [`python-api`](python-api)                                                                                           | [libhoney-py](https://docs.honeycomb.io/sdk/python/)                                                                           | A TODO API written and instrumented using Python (Flask).                                                                                                                                                 |
| [`beeline-python`](https://github.com/honeycombio/beeline-python/tree/main/examples)                                 | [Python Beeline](https://docs.honeycomb.io/getting-data-in/beeline/python/)                                                    | Examples of using the Beeline Python SDK                                                                                                                                                                  |
| [`beeline-ruby`](https://github.com/honeycombio/beeline-ruby/tree/main/examples)                                     | [Ruby Beeline](https://docs.honeycomb.io/getting-data-in/beeline/ruby/)                                                        | Examples of using the Beeline Ruby SDK                                                                                                                                                                    |
| [`ruby-wiki-tracing`](ruby-wiki-tracing)                                                                             | [Manual Tracing](https://docs.honeycomb.io/working-with-data/tracing/send-trace-data/#manual-tracing) with Ruby                | A simple wiki (Ruby) manually instrumented for tracing.                                                                                                                                                   |


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
