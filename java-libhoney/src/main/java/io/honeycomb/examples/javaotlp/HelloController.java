package io.honeycomb.examples.javaotlp;

import io.honeycomb.libhoney.*;
import org.springframework.web.bind.annotation.*;

@RestController
public class HelloController {

    @RequestMapping("/")
    public String index() {
        final HoneyClient honeyClient = Honey.getHoneyClient();
        final Event event = honeyClient.createEvent();
        event.addField("request.name", "Hello");
        try {
            return "Hello world!";
        } finally {
            event.sendPresampled();
        }
    }

}
