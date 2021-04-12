package io.honeycomb.examples.javaotlp;

import io.opentelemetry.api.baggage.*;
import io.opentelemetry.context.*;
import io.opentelemetry.sdk.trace.*;

public class BaggageSpanProcessor implements SpanProcessor {
    @Override
    public void onStart(Context parentContext, ReadWriteSpan span) {
        Baggage.fromContext(parentContext)
                .forEach((s, baggageEntry) -> span.setAttribute(s, baggageEntry.getValue()));
    }

    @Override
    public boolean isStartRequired() {
        return true;
    }

    @Override
    public void onEnd(ReadableSpan span) {
    }

    @Override
    public boolean isEndRequired() {
        return false;
    }
}
