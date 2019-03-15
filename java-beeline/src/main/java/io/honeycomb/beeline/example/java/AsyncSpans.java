package io.honeycomb.beeline.example.java;

import io.honeycomb.beeline.tracing.Beeline;
import io.honeycomb.beeline.tracing.Span;
import io.honeycomb.beeline.tracing.SpanBuilderFactory;
import io.honeycomb.beeline.tracing.SpanPostProcessor;
import io.honeycomb.beeline.tracing.Tracer;
import io.honeycomb.beeline.tracing.Tracing;
import io.honeycomb.beeline.tracing.sampling.Sampling;
import io.honeycomb.libhoney.HoneyClient;
import io.honeycomb.libhoney.LibHoney;

import java.util.concurrent.Callable;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;

import static java.util.concurrent.Executors.newSingleThreadExecutor;

/**
 * This demonstrates instrumenting multi-threaded code using the Beeline and its collaborator classes.
 * <p>
 * When using the Beeline in a Spring Boot application all the setup code would be handled by the Beeline
 * AutoConfiguration.
 */
// @formatter:off
@SuppressWarnings("ALL")
public class AsyncSpans {
    private static final String WRITE_KEY = "test-write-key";
    private static final String DATASET   = "test-dataset";

    private static final Beeline beeline;
    private static final HoneyClient client;

    static {
        client                          = LibHoney.create(LibHoney.options().setDataset(DATASET).setWriteKey(WRITE_KEY).build());
        SpanPostProcessor postProcessor = Tracing.createSpanProcessor(client, Sampling.alwaysSampler());
        SpanBuilderFactory factory      = Tracing.createSpanBuilderFactory(postProcessor, Sampling.alwaysSampler());
        Tracer tracer                   = Tracing.createTracer(factory);
        beeline                         = Tracing.createBeeline(tracer, factory);
    }

    private HttpClient httpClient = new HttpClient();
    private String requestUrl = "https://example.com";

    public static void main(String... args) throws Exception {
        AsyncSpans examples = new AsyncSpans();
        Span rootSpan = beeline.getSpanBuilderFactory()
            .createBuilder()
            .setSpanName("base-span")
            .setServiceName("async-service")
            .build();

        try(Span tracedSpan = beeline.getTracer().startTrace(rootSpan)) {
            examples.run();
        } finally{
            // close to flush events and release its thread pool
            client.close();
        }
    }

    /**
     * Runs the various examples. They all create children to the root span created in the main method.
     */
    void run() throws Exception {
        futureChain();

        singleLambda();

        multipleSpansInAPipeline();

        submitCallable();

        submitRunnable();
    }

    /**
     * This demonstrates how one could create a "detached" Span to cover a CompletableFuture chain.
     * <p>
     * We cannot simply use {@link Beeline#getActiveSpan()}, which uses thread-locals, within the CompletableFuture,
     * because runAsync will run the task on a different thread
     */
    private void futureChain() {
        // create a child span that is not lonked to the Beeline's thread local context.
        Span httpServiceSpan = beeline.getTracer().startDetachedChildSpan("http-call-futureChain");

        httpServiceSpan.addField("http-request-url", requestUrl);

        CompletableFuture
            // start measuring duration when execution of the task starts
            .runAsync(httpServiceSpan::markStart)

            // perform operations
            .thenApply((v) -> httpClient.get(requestUrl))
            .thenApply(this::convertResponse)
            .thenApply(this::handleResponse)

            // add any error that might have occured in the above steps
            .exceptionally(e -> httpServiceSpan.addField("error-message", e.getMessage()))

            // close and submit span to honeycomb
            .thenRun(httpServiceSpan::close);
    }

    /**
     * This is functionally the same as the "futureChain" example except that it we turned the chain into a single
     * lambda.
     */
    private void singleLambda() {
        Span httpServiceSpan = beeline.getTracer().startDetachedChildSpan("http-call-lambda");

        httpServiceSpan.addField("http-request-url", requestUrl);
        CompletableFuture
            .runAsync(() -> {
                httpServiceSpan.markStart();
                try {
                    String rawResponse = httpClient.get(requestUrl);
                    String response = convertResponse(rawResponse);
                    handleResponse(response);
                } catch (Throwable e) {
                    httpServiceSpan.addField("error-message", e.getMessage());
                } finally {
                    httpServiceSpan.close();
                }
            });
    }

    /**
     * In the "futureChain" and "singleLambda" examples we create a detached child Span, because we run the task on a
     * different thread, However, we can also propagate the thread-local context of the Beeline to the new thread,
     * and use it inside the task, by wrapping the task Runnable.
     */
    private void submitRunnable() throws ExecutionException, InterruptedException {
        Runnable task = () -> {
            beeline.getActiveSpan().addField("http-request-url", requestUrl);
            httpClient.get(requestUrl);
        };

        Runnable tracedTask = beeline.getTracer().traceRunnable("http-call-runnable", task);

        ExecutorService executorService = newSingleThreadExecutor();
        executorService.submit(tracedTask).get();
        executorService.shutdownNow();
    }

    /**
     * Same as the "submitRunnable" example, except with a Callable. This simply demonstrates that different "functional
     * interfaces" that are supported.
     */
    private void submitCallable() throws ExecutionException, InterruptedException {
        Callable<String> task = () -> {
            beeline.getActiveSpan().addField("http-request-url", requestUrl);
            return httpClient.get(requestUrl);
        };

        Callable<String> tracedTask = beeline.getTracer().traceCallable("http-call-callable", task);

        ExecutorService executorService = newSingleThreadExecutor();
        executorService.submit(tracedTask).get();
        executorService.shutdownNow();
    }

    /**
     * This further demonstrates that Supplier and Function can also be wrapped.
     * In contrast to the "futureChain" example, this would create 3 Span children for each wrapped lambda.
     */
    private void multipleSpansInAPipeline() {
        CompletableFuture
            .supplyAsync(
                beeline.getTracer().traceSupplier("http-call-first", () -> httpClient.get(requestUrl)))
            .thenApplyAsync(
                beeline.getTracer().traceFunction("http-call-conversion", this::convertResponse))
            .thenApplyAsync(
                beeline.getTracer().traceFunction("http-call-handle", this::handleResponse));
    }

    // placeholders to make the examples above look more "realistic"
    private String convertResponse(String s) {
        return "";
    }

    private Object handleResponse(String response) {
        return "";
    }

    public static class HttpClient {
        public String get(String requestUrl) {
            return "mock response";
        }
    }
}

