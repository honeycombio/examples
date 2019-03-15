package io.honeycomb.beeline.example.spring;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;

@SuppressWarnings("ALL")
@Controller
public class DemoForwardingController {
    /**
     * Forward call to another endpoint, showing nesting of Spans when performing a "forward dispatch".
     *
     * @param redirectPath to use when forwarding.
     * @return forward.
     */
    @RequestMapping("/forward-to/{redirectPath}")
    String forwardingEndpoint(@PathVariable String redirectPath) {
        return "forward:/" + redirectPath;
    }
}
