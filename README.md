# Honeycomb Instrumentation Examples

[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/examples?color=success)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)

This repository is an index of various examples showing how to send data to Honeycomb. Most of the Honeycomb tools and SDKs have living examples in their respective repositories, listed below.

| Examples                                                                                                             | Keywords                                                            | Description                                                                                                                                                                                               |
| -------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`honeycomb-opentelemetry-java`](https://github.com/honeycombio/honeycomb-opentelemetry-java/tree/main/examples)     | OpenTelemetry, traces, auto-instrumentation                         | Spring Boot applications instrumented using the [Honeycomb OpenTelemetry Java SDK Distribution](https://docs.honeycomb.io/getting-data-in/opentelemetry/java-distro/)                                     |
| [`honeycomb-opentelemetry-dotnet`](https://github.com/honeycombio/honeycomb-opentelemetry-dotnet/tree/main/examples) | OpenTelemetry, traces, auto-instrumentation                         | .NET applications instrumented using the [Honeycomb OpenTelemetry .NET SDK Distribution](https://docs.honeycomb.io/getting-data-in/opentelemetry/dotnet-distro/)                                          |
| [`honeycomb-opentelemetry-node`](https://github.com/honeycombio/honeycomb-opentelemetry-node/tree/main/examples)     | OpenTelemetry, traces, auto-instrumentation                         | Node.js applications instrumented using the [Honeycomb OpenTelemetry Node.js SDK Distribution](https://docs.honeycomb.io/getting-data-in/opentelemetry/node-distro/)                                      |
| [`honeycomb-opentelemetry-python`](https://github.com/honeycombio/honeycomb-opentelemetry-python/tree/main/examples) | OpenTelemetry, traces, auto-instrumentation                         | Python applications instrumented using the [Honeycomb OpenTelemetry Python SDK Distribution](https://docs.honeycomb.io/getting-data-in/opentelemetry/python-distro/)                                      |
| [`honeycomb-opentelemetry-go`](https://github.com/honeycombio/honeycomb-opentelemetry-go/tree/main/examples)         | OpenTelemetry, traces                                               | Go applications instrumented using the [Honeycomb OpenTelemetry Go SDK Distribution](https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/)                                                  |
| [`beeline-nodejs`](https://github.com/honeycombio/beeline-nodejs/tree/main/examples/node-tracing)                    | Beeline, traces                                                     | A simple webapp instrumented with the [NodeJS Beeline](https://docs.honeycomb.io/getting-data-in/beeline/nodejs/).                                                                                        |
| [`beeline-python`](https://github.com/honeycombio/beeline-python/tree/main/examples)                                 | Beeline, traces, auto-instrumentation, flask                        | Examples of using the [Beeline Python SDK](https://docs.honeycomb.io/getting-data-in/beeline/python/)                                                                                                     |
| [`beeline-ruby`](https://github.com/honeycombio/beeline-ruby/tree/main/examples)                                     | Beeline, traces, auto-instrumentation, rails, sinatra, rack, sequel | Examples of using the [Beeline Ruby SDK](https://docs.honeycomb.io/getting-data-in/beeline/ruby/)                                                                                                         |
| [`Example Greeting Service`](https://github.com/honeycombio/example-greeting-service)                                | Microservices, traces, OpenTelemetry, Beelines                      | A ridiculously over-engineered microservice application implemented in several languages. The services are instrumented to send telemetry to Honeycomb, with examples of Beelines and OpenTelemetry SDKs. |
| [`honeytail`](https://github.com/honeycombio/honeytail/tree/main/examples)                                           | Logs, haproxy, mysql, nginx                                         | Examples of using [Honeytail](https://docs.honeycomb.io/getting-data-in/logs/honeytail/) to ingest structured log files.                                                                                  |
| [`libhoney-go`](https://github.com/honeycombio/libhoney-go/tree/main/examples)                                       | Structured events, manual traces                                    | Examples of using the low-level [Libhoney Go SDK](https://docs.honeycomb.io/getting-data-in/libhoney/go/) to send structured events to Honeycomb, including manual tracing.                               |
| [`libhoney-java`](https://github.com/honeycombio/libhoney-java/tree/main/examples)                                   | Structures events                                                   | Examples of using the low-level [Libhoney Java SDK](https://docs.honeycomb.io/getting-data-in/libhoney/java/) to send structured events to Honeycomb                                                      |
| [`libhoney-py`](https://github.com/honeycombio/libhoney-py/tree/main/examples)                                       | Structure events, manual traces                                     | Examples of using the low-level [Libhoney Python SDK](https://docs.honeycomb.io/getting-data-in/libhoney/python/) to send structured events to Honeycomb, including manual tracing.                       |
| [`libhoney-rb`](https://github.com/honeycombio/libhoney-rb/tree/main/examples)                                       | Structured events, manual traces                                    | Examples of using the low-level [Libhoney Ruby SDK](https://docs.honeycomb.io/getting-data-in/libhoney/ruby/) to send structured events to Honeycomb, including manual tracing.                           |
| [`buildevents`](https://github.com/honeycombio/buildevents)                                                          | Manual traces, Libhoney Go                                          | A small binary to help create traces out of CI builds, visualizing each step and command as spans within the trace.                                                                                       |
