package io.honeycomb.beeline.example.spring.services;

import io.honeycomb.beeline.spring.beans.aspects.ChildSpan;
import io.honeycomb.beeline.spring.beans.aspects.SpanField;
import io.honeycomb.beeline.tracing.Beeline;
import io.honeycomb.beeline.tracing.Tracer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.concurrent.CompletableFuture;

@SuppressWarnings("ALL")
@Component
public class AspectAnnotatedComponent {
    @Autowired
    Beeline beeline;

    /**
     * Declaratively generate a child span when the method is invoked. The Span name is inferred from the method name.
     * Additionally, the annotation requests that the return value is added as a field.
     *
     * @param captains to pick from.
     * @return one of the captains.
     */
    @ChildSpan(addResult = true)
    public String pickARandomCaptain(String... captains) {
        addLatency(100);
        return captains[(int) (Math.random() * captains.length)];
    }

    /**
     * Demonstrates that an explicit Span name can be set, and that the arguments can declaratively be captured.
     *
     * @param title       to use a prefix to the name.
     * @param captainName to use.
     * @return title + name.
     */
    @ChildSpan("Add-Title-To-Name")
    public String prefixWithTitle(@SpanField("captain-title") String title, @SpanField String captainName) {
        addLatency(100);
        return title + " " + captainName;
    }

    /**
     * Demonstrates propagation to other threads and how async spans can outlive the root.
     * There are various calls to Thread.sleep to make the span durations more "realistic".
     */
    @ChildSpan
    public void runAsyncTasks() {
        addLatency(100);
        final Tracer tracer = beeline.getTracer();

        CompletableFuture.runAsync(tracer.traceRunnable("task-1", () -> {
            addLatency(300);
            beeline.getActiveSpan().addField("task", 1);
        }));
        CompletableFuture.runAsync(tracer.traceRunnable("task-2", () -> {
            addLatency(200);
            beeline.getActiveSpan().addField("task", 2);
        }));

        addLatency(300);
        CompletableFuture
            .runAsync(tracer.traceRunnable("task-3-1", () -> {
                addLatency(300);
                beeline.getActiveSpan().addField("task", 3.1);
            }))
            .thenRun(tracer.traceRunnable("task-3-2", () -> {
                addLatency(300);
                beeline.getActiveSpan().addField("task", 3.2);
            }));
        addLatency(100);
    }

    /**
     * Creates some latency to make Span times more realistic.
     */
    private void addLatency(int waitMultiplyer) {
        try {
            Thread.sleep((long) (Math.random() * waitMultiplyer));
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
}
