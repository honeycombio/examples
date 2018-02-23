# Honeycomb Instrumentation Examples

The full power of Honeycomb is unlocked by proper instrumentation, especially at
the code level. The examples in this respository are intended to help guide you
along the way to unlocking this power for yourself by understand how
instrumentation is meant to be done.

Most of the top level directories in this repository correspond to an example
which demonstrates instrumentation using Honeycomb and sample queries to help
you along. The current examples are:

| Directory | Description |
| --- | --- |
| `golang-webapp` | A two-tier web application (Go+MySQL) which is a Twitter clone. |
| `honeytail-dockerd` | Using [Honeytail]()'s `keyval` parser to ingest the structured logs of the [Docker]() container engine daemon. |
| `honeytail-nginx` | Using [Honeytail]()'s `nginx` parser to ingest [Nginx]() access logs from an instance acting as a reverse proxy. |
| `python-api` | A TODO API written and instrumented using Python (Flask). |

## Proposed Examples

The following have been proposed but not implemented:

| Directory | Description |
| --- | --- |
| `java-api` | A TODO API written and instrumented using Java. |
| `golang-api` | A TODO API written and instrumented using Golang. |
| `ruby-api` | A TODO API written and instrumented using Ruby. |
| `javascript-api` | A TODO API written and instrumented using JavaScript. |
| `sidekiq` | Observing behavior of the background job runner Sidekiq. |
| `honeytail-apache` | Ingesting Apache access logs using Honeytail. |
| `honeytail-haproxy` | Ingesting HAProxy access logs using Honeytail. |
| `logstash` | Using the Honeycomb Logstash plugin to send parsed events to Honeycomb. |
| `fluentd` | Using the Honeycomb Fluentd plugin to send parsed events to Honeycomb. |

Let us know if there is something specific you'd like to see by [filing an
issue](https://github.com/honeycombio/examples/issues/new).
