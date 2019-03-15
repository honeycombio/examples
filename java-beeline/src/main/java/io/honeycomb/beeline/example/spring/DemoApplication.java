package io.honeycomb.beeline.example.spring;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.common.Slf4jNotifier;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.web.client.RestTemplate;

import javax.annotation.PreDestroy;

import static com.github.tomakehurst.wiremock.client.WireMock.aResponse;
import static com.github.tomakehurst.wiremock.client.WireMock.post;
import static com.github.tomakehurst.wiremock.client.WireMock.urlPathMatching;

@SuppressWarnings("ALL")
@SpringBootApplication
public class DemoApplication {

    public static void main(String... args) {
        SpringApplication.run(DemoApplication.class, args);
    }

    ///////////////////////////////////////////////////////////////////////////
    // Create an instrumented rest template.
    // The Beeline AutoConfiguration wires itself into the RestTemplateBuilder.
    ///////////////////////////////////////////////////////////////////////////
    @Bean
    public RestTemplate restTemplate(RestTemplateBuilder builder) {
        return builder.build();
    }

    ///////////////////////////////////////////////////////////////////////////
    // Start up a mock server (running on localhost:8089) that stubs Honeycomb's Events API.
    // We target it with the "honeycomb.beeline.api-host" property.
    ///////////////////////////////////////////////////////////////////////////
    @Autowired
    WireMockServer mockServer;

    @PreDestroy
    public void closeMock() {
        mockServer.stop();
    }

    @Bean(destroyMethod = "")
    public WireMockServer honeycombMock() {
        WireMockServer honeycombMock = new WireMockServer(WireMockConfiguration.options()
            .port(8089)
            .notifier(new Slf4jNotifier(false)));
        honeycombMock.start();
        stubEventsEndpoint(honeycombMock);
        return honeycombMock;
    }

    private static void stubEventsEndpoint(WireMockServer honeycombMock) {
        honeycombMock.stubFor(post(urlPathMatching("/1/batch/.*"))
            .willReturn(aResponse().withBody("[{\"status\": 202},{\"status\": 202}, {\"status\": 202}, {\"status\": 202}, {\"status\": 202}, {\"status\": 202}, {\"status\": 202}]")));
    }
}

