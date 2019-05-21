# node-tracing-webapp

## Overview

This example illustrates a simple comment-wall application using Honeycomb's [Beeline for Node](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/).

It contains examples of:

- Baseline Beeline usage in an Express app
- Capture of custom metadata on Beeline-generated events
- Definition of custom spans to augment traces
- In-app tracing (across functions within the same service)
- Distributed tracing (across services)

It requires Node 8+ and has intentionally included the `node_modules` directory in hopes of simplifying dependencies.

## Usage:

Find your Honeycomb write key at https://ui.honeycomb.io/account, then:

## Running the main application

Run our sample `wall` service with:

```bash
# Will run on port 8080
$ HONEYCOMB_API_KEY=foobarbaz node ./wall.js
```

### Interacting with your application

You may either use the web UI [`http://localhost:8080`](http://localhost:8080) to read and write messages:

![index](/images/index.png) | ![new message](/images/message.png)
:-------------------------:|:-------------------------:
View contents of wall | Write new message on wall

Or `curl` the contents of your wall directly:

```bash
# Fetch the contents of your wall
curl localhost:8080
```

```bash
# Write a new message to your wall
curl localhost:8080 -d "message=i'm #tracing with @honeycombio"
```

### Running the analysis service

This webapp may call out to a second service, `analysis`, with:

```bash
# Hard-coded to run on port 8088
$ HONEYCOMB_API_KEY=foobarbaz node ./analysis.js
```

But you won't be interacting with it directly; the `wall` service will simply ping `localhost:8088` in hopes of the `analysis` service being alive.
