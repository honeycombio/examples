package io.honeycomb.beeline.example.java;

import io.honeycomb.beeline.tracing.Beeline;
import io.honeycomb.beeline.tracing.Span;
import io.honeycomb.beeline.tracing.SpanBuilderFactory;
import io.honeycomb.beeline.tracing.SpanPostProcessor;
import io.honeycomb.beeline.tracing.Tracer;
import io.honeycomb.beeline.tracing.Tracing;
import io.honeycomb.beeline.tracing.propagation.HttpHeaderV1PropagationCodec;
import io.honeycomb.beeline.tracing.propagation.Propagation;
import io.honeycomb.beeline.tracing.propagation.PropagationContext;
import io.honeycomb.beeline.tracing.sampling.Sampling;
import io.honeycomb.libhoney.HoneyClient;
import io.honeycomb.libhoney.LibHoney;

/**
 * This shows how one could manually propagate Spans. This is in contrast with the examples in {@link TracerSpans} or
 * {@link AsyncSpans}, which use the convenience of the Beeline and its collaborator classes to do so.
 * <p>
 * When using the Beeline in a Spring Boot application all the setup code would be handled by the Beeline
 * AutoConfiguration.
 */
// @formatter:off
@SuppressWarnings("ALL")
public class TracerSpans {
    private static final String WRITE_KEY = "test-write-key";
    private static final String DATASET   = "test-dataset";

    private static final HoneyClient client;
    private static final Beeline beeline;

    static {
        client                          = LibHoney.create(LibHoney.options().setDataset(DATASET).setWriteKey(WRITE_KEY).build());
        SpanPostProcessor postProcessor = Tracing.createSpanProcessor(client, Sampling.alwaysSampler());
        SpanBuilderFactory factory      = Tracing.createSpanBuilderFactory(postProcessor, Sampling.alwaysSampler());
        Tracer tracer                   = Tracing.createTracer(factory);
        beeline                         = Tracing.createBeeline(tracer, factory);
    }

    public static void main(String... args) {
        TracerSpans example = new TracerSpans();
        try {
            HttpRequest request = new HttpRequest();
            startTrace(request);
            example.acceptRequest(request);
        } finally {
            beeline.getTracer().endTrace();
            client.close(); // close to flush events and release its thread pool
        }
    }

    private DatabaseService db = new DatabaseService();

    private static void startTrace(HttpRequest request) {
        String headerValue = request.getHeader(HttpHeaderV1PropagationCodec.HONEYCOMB_TRACE_HEADER);
        PropagationContext context = Propagation.honeycombHeaderV1().decode(headerValue);

        Span rootSpan = beeline.getSpanBuilderFactory().createBuilder()
            .setSpanName("get-customer-data")
            .setServiceName("customer-db-traced")
            .setParentContext(context)
            .build();
        beeline.getTracer().startTrace(rootSpan);
    }

    // @RequestMapping
    public void acceptRequest(HttpRequest request) {
        Span span = beeline.getActiveSpan();
        try {
            db.queryDb(request.getParameter("customer-id"));
            span.addField("result", "OK");
        } catch (Exception e) {
            span.addField("result", "Bad Request")
                .addField("exception-message", e.getMessage());
        }
    }

    public static class DatabaseService {
        public void queryDb(String id) {
            try (Span childSpan = beeline.startChildSpan("customer-db-query")) {
                String data = getCustomerDataById(id);
                childSpan.addField("customer-data", data);
            }

        }

        public String getCustomerDataById(String id) {
            return "customer-0123";
        }

    }

    private static class HttpRequest {
        public String getHeader(String key) {
            return "mockHeader";
        }

        public String getRequestURI() {
            return "mockHeader";
        }

        public String getParameter(String key) {
            return "mockParameter";
        }
    }
}
