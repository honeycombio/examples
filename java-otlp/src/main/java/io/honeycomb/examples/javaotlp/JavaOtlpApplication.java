package io.honeycomb.examples.javaotlp;

import java.util.Arrays;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.common.AttributeKey;
import io.opentelemetry.api.common.Attributes;
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
			.setEndpoint("https://api.honeycomb.io:443")
			.addHeader("x-honeycomb-team", "<YOUR-APIKEY>")
			.addHeader("x-honeycomb-dataset", "<YOUR-DATASET>")
			.build();

		// Configure the OpenTelemtry SDK with the exporter
		OpenTelemetrySdk.builder()
			.setTracerProvider(
				SdkTracerProvider.builder()
					.setResource(Resource.create(Attributes.of(AttributeKey.stringKey("service.name"), "java-otlp")))
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
