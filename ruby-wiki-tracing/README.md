# ruby-wiki-tracing

This example illustrates a simple wiki application instrumented with the bare minimum necessary to utilize Honeycomb's tracing functionality.

## What We'll Do in This Example

We'll instrument a simple application for tracing by following a few general steps:

1. Set a top-level `traceId` at the origin of the request and set it on the request context. Generate a root span indicated by omitting a `parentId` field.
2. To represent a unit of work within a trace as a span, add code to generate a span ID and capture the start time. At the **call site** of the unit of work, pass down a new request context with the newly-generated span ID as the `parentId`. Upon work completion, send the span with a calculated `durationMs`.
3. Rinse and repeat.

**Note**: The [Honeycomb Beeline for Ruby](https://github.com/honeycombio/beeline-ruby) handles all of this propagation magic for you :)

## Usage

Find your Honeycomb write key at https://ui.honeycomb.io/account, then run:

```bash
$ HONEYCOMB_WRITEKEY=foobarbaz ruby wiki.rb
```

And load [`http://localhost:4567/view/MyFirstWikiPage`](http://localhost:4567/view/MyFirstWikiPage) to create (then view) your first wiki page.

Methods within the simple wiki application have been instrumented with tracing-like calls, with tracing identifiers propagated via thread locals.

![Honeycomb initial view](images/honeycomb.png?raw=true "Explore traces with Honeycomb")

And click into the graph to inspect individual traces:

![Honeycomb trace view](images/trace.png?raw=true "Classic tracing waterfall view in Honeycomb")

## Tips for Instrumenting Your Own Service

- For a given span (e.g. `"loadPage"`), remember that the span definition lives in its parent, and the instrumentation is around the **call site** of `loadPage`:
    ```ruby
    with_span("load_page") do
      # sets the appropriate "parent id" within the scope of the block
      load_page(title)
      # span is sent to Honeycomb upon completion of the block
    end
    ```
- If emitting Honeycomb events or structured logs, make sure that the **start** time gets used as the canonical timestamp, not the time at event emission.
- Remember, the root span should **not** have a `parentId`.
- Don't forget to add some metadata of your own! It's helpful to identify metadata in the surrounding code that might be interesting when debugging your application.
- Check out our [Honeycomb Beeline for Ruby](https://github.com/honeycombio/beeline-ruby) to get this context propagation for free! 

## A Note on Code Style

The purpose of this example is to illustrate the **bare minimum necessary** to propagate and set identifiers to enable tracing on an application for consumption by Honeycomb, illustrating the steps described in the top section of this README.

Are you grossed out by the code? Fortunately, you don't have to DIY: check out the [Honeycomb Beeline for Ruby](https://github.com/honeycombio/beeline-ruby)!
