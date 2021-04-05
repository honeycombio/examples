package io.honeycomb.examples.javaotlp;

import java.util.Arrays;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.baggage.propagation.*;
import io.opentelemetry.api.common.AttributeKey;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.trace.propagation.*;
import io.opentelemetry.context.propagation.*;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.resources.Resource;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor;

import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;

@SpringBootApplication
public class JavaOtlpApplication {

    public static void main(String[] args) {
        SpringApplication.run(JavaOtlpApplication.class, args);

        // Create OTLP exporter that sends trace data to Honeycomb
        OtlpGrpcSpanExporter exporter = OtlpGrpcSpanExporter.builder()
            .setEndpoint(System.getenv("HONEYCOMB_API_ENDPOINT"))
            .addHeader("x-honeycomb-team", System.getenv("HONEYCOMB_API_KEY"))
            .addHeader("x-honeycomb-dataset", System.getenv("HONEYCOMB_DATASET"))
            .build();
//
        // Configure the OpenTelemtry SDK with the exporter
        ContextPropagators propagators = ContextPropagators.create(
            TextMapPropagator.composite(
                W3CTraceContextPropagator.getInstance(),
                W3CBaggagePropagator.getInstance()));
        OpenTelemetrySdk.builder()
            .setPropagators(propagators)
            .setTracerProvider(
                SdkTracerProvider.builder()
                    .setResource(Resource.create(Attributes.of(AttributeKey.stringKey("service.name"), "manual-instrumentation")))
                    .addSpanProcessor(new BaggageSpanProcessor())
                    .addSpanProcessor(BatchSpanProcessor.builder(exporter).build())
                    .build())
            .buildAndRegisterGlobal();
    }

    @Bean
    public CommandLineRunner commandLineRunner(ApplicationContext ctx) {
        return args -> {

            String[] beanNames = ctx.getBeanDefinitionNames();
            Arrays.sort(beanNames);
            for (String beanName : beanNames) {
                System.out.println(beanName);
            }

        };
    }

}
