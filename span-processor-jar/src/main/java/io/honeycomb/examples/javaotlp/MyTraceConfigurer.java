package io.honeycomb.examples.javaotlp;

import io.opentelemetry.sdk.autoconfigure.spi.*;
import io.opentelemetry.sdk.trace.*;


public class MyTraceConfigurer implements SdkTracerProviderConfigurer {
    @Override
    public void configure(SdkTracerProviderBuilder tracerProvider) {
        tracerProvider.addSpanProcessor(new BaggageSpanProcessor());
    }
}
