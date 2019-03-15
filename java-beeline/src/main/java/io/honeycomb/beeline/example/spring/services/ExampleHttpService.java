package io.honeycomb.beeline.example.spring.services;

import io.honeycomb.beeline.tracing.Beeline;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

@SuppressWarnings("ALL")
@Component
public class ExampleHttpService {
    private static final Logger LOG = LoggerFactory.getLogger(ExampleHttpService.class);

    private static final String EXAMPLE_COM_URL = "http://www.example.com";

    @Autowired
    Beeline beeline;

    @Autowired
    RestTemplate client;

    public String getExampleContent() {
        String response = client.getForObject(EXAMPLE_COM_URL, String.class);
        LOG.info("Received response from example.com: {}", response.substring(0, 100));
        beeline.getActiveSpan().addField("example-response-substring", response.substring(0, 100));
        return response;
    }
}
