package io.honeycomb.examples.javaotlp;

import java.util.Arrays;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
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

		System.setProperty("otel.resource.attributes", "service.name=java-otlp");

		// Create OTLP exporter that sends trace data to Honeycomb
		OtlpGrpcSpanExporter exporter = OtlpGrpcSpanExporter.builder()
			.setEndpoint("api.honeycomb.io:443")
			.setUseTls(true)
			.addHeader("x-honeycomb-team", "")
			.addHeader("x-honeycomb-dataset", "test-otlp")
			.build();

		// Configure the OpenTelemtry SDK with the exporter
		OpenTelemetrySdk.getGlobalTracerManagement()
			.addSpanProcessor(BatchSpanProcessor.builder(exporter).build());
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
