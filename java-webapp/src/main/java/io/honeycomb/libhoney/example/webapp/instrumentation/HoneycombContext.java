package io.honeycomb.libhoney.example.webapp.instrumentation;

import io.honeycomb.libhoney.Event;
import io.honeycomb.libhoney.HoneyClient;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.context.annotation.RequestScope;

import javax.annotation.PostConstruct;
import java.util.HashMap;
import java.util.Map;

@RequestScope
@Component
public class HoneycombContext {
    @Autowired
    private HoneyClient honeyClient;

    private Map<String, Object> contextualData;
    private Event event;

    @PostConstruct
    public void setup() {
        this.event = honeyClient.createEvent();
        this.contextualData = new HashMap<>();
    }

    public Map<String, Object> getContextualData() {
        return contextualData;
    }

    public Event getEvent() {
        return event;
    }
}
