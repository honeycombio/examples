package io.honeycomb.libhoney.example.webapp.configuration;

import io.honeycomb.libhoney.HoneyClient;
import io.honeycomb.libhoney.LibHoney;
import io.honeycomb.libhoney.ResponseObserver;
import io.honeycomb.libhoney.example.webapp.LoggingResponseObserver;
import io.honeycomb.libhoney.example.webapp.instrumentation.HoneycombHandlerInterceptor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

import java.net.URI;
import java.net.URISyntaxException;

@Configuration
@EnableConfigurationProperties(HoneycombProperties.class)
public class ApplicationConfiguration extends WebMvcConfigurerAdapter {
    @Autowired
    private HoneycombHandlerInterceptor honeycombHandlerInterceptor;

    @Bean
    public HoneyClient honeyClient(final HoneycombProperties honeycombProperties) {
        try {
            final HoneyClient honeyClient = LibHoney.create(LibHoney.options()
                .setDataset(honeycombProperties.getDataset())
                .setWriteKey(honeycombProperties.getWriteKey())
                .setSampleRate(honeycombProperties.getSampleRate())
                .setApiHost(honeycombProperties.getApiHost() == null ? null : new URI(honeycombProperties.getApiHost()))
                .setGlobalFields(honeycombProperties.getGlobalFields())
                .build());

            honeyClient.addResponseObserver(responseObserver());

            return honeyClient;
        } catch (final URISyntaxException ex) {
            throw new IllegalArgumentException("Error when creating URI from configured APIHost string", ex);
        }
    }

    @Bean
    public ResponseObserver responseObserver() {
        return new LoggingResponseObserver();
    }

    @Override
    public void addInterceptors(final InterceptorRegistry registry) {
        registry.addInterceptor(honeycombHandlerInterceptor);
    }
}
