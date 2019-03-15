package io.honeycomb.beeline.example.java;

import io.honeycomb.beeline.tracing.Span;
import io.honeycomb.beeline.tracing.SpanBuilderFactory;
import io.honeycomb.beeline.tracing.SpanPostProcessor;
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
public class ManualSpans {
    private static final String WRITE_KEY = "test-write-key";
    private static final String DATASET   = "test-dataset";

    private static final SpanBuilderFactory factory;
    private static final HoneyClient client;
    static {
        client                          = LibHoney.create(LibHoney.options().setDataset(DATASET).setWriteKey(WRITE_KEY).build());
        SpanPostProcessor postProcessor = Tracing.createSpanProcessor(client, Sampling.alwaysSampler());
        factory                         = Tracing.createSpanBuilderFactory(postProcessor, Sampling.alwaysSampler());
    }

    public static void main(String... args){
        ManualSpans example = new ManualSpans();
        HttpRequest request = new HttpRequest();
        try(Span rootSpan = startTrace(request)) {
            example.acceptRequest(request, rootSpan);
        } finally {
            client.close(); // close to flush events and release its thread pool
        }
    }

    private static Span startTrace(HttpRequest request) {
        // extract propagated trace from header, if it exists
        String headerValue = request.getHeader(HttpHeaderV1PropagationCodec.HONEYCOMB_TRACE_HEADER);
        PropagationContext context = Propagation.honeycombHeaderV1().decode(headerValue);

        // set up root span with some context data
        return factory.createBuilder()
            .setSpanName("get-customer-data")
            .setServiceName("customer-db-manual")
            .setParentContext(context)
            .addField("request-uri", request.getRequestURI())
            .addField("customer-id", request.getParameter("customer-id"))
            .build();
    }

    private DatabaseService db = new DatabaseService();

    // @RequestMapping
    public void acceptRequest(HttpRequest request, Span span) {
        try {
            db.queryDb(request.getParameter("customer_id"), span);
            span.addField("result", "OK");
        } catch (Exception e) {
            span.addField("result", "Bad Request")
                    .addField("exception-message", e.getMessage());
        }
    }

    public static class DatabaseService {
        public void queryDb(String id, Span parentSpan) {
            SpanBuilderFactory.SpanBuilder spanBuilder = factory.createBuilderFromParent(parentSpan).setSpanName("customer-db-query");
            try (Span childSpan = spanBuilder.build()) {
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
