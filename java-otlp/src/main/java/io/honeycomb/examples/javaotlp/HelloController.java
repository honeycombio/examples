package io.honeycomb.examples.javaotlp;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestMapping;

@RestController
public class HelloController {

	@RequestMapping("/")
	public String index() {
		Span span = GlobalOpenTelemetry.getTracer("java-otlp")
			.spanBuilder("hello")
			.startSpan();
		try {
			return "Hello world!";
		} finally {
			span.end();
		}
	}

}
