package com.example.javaagent;

import io.opentelemetry.sdk.autoconfigure.spi.SdkTracerProviderConfigurer;
import io.opentelemetry.sdk.trace.SdkTracerProviderBuilder;

public class DemoSdkTracerProviderConfigurer implements SdkTracerProviderConfigurer {
  @Override
  public void configure(SdkTracerProviderBuilder tracerProvider) {
    tracerProvider
        .addSpanProcessor(new BaggageSpanProcessor());
  }
}
