package io.honeycomb.beeline.example.spring;

import io.honeycomb.beeline.example.spring.data.CustomerData;
import io.honeycomb.beeline.example.spring.services.AspectAnnotatedComponent;
import io.honeycomb.beeline.example.spring.services.ExampleHttpService;
import io.honeycomb.beeline.spring.beans.aspects.ChildSpan;
import io.honeycomb.beeline.tracing.Beeline;
import io.honeycomb.beeline.tracing.Span;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@SuppressWarnings("ALL")
@RestController
public class DemoRestController {
    @Autowired
    Beeline beeline;

    @Autowired
    AspectAnnotatedComponent component;

    @Autowired
    ExampleHttpService exampleHttpService;

    /**
     * Make HTTP client call to example.com, demonstrating server + client spans.
     *
     * @return the HTML response.
     */
    @RequestMapping(value = "/get-example", produces = MediaType.TEXT_HTML_VALUE)
    String getExamplePage() {
        return exampleHttpService.getExampleContent();
    }

    /**
     * Concatenate query parameters (or use the default value).
     * Shows how the "request.query" field is populated.
     *
     * @param firstWord  to concatenate.
     * @param secondWord to concatenate.
     * @return the concatenation.
     */
    @RequestMapping("/concatenate-words")
    String concatenateWords(@RequestParam(name = "first-word", defaultValue = "Hello") String firstWord,
                            @RequestParam(name = "second-word", defaultValue = "World") String secondWord) {
        String result = firstWord.toUpperCase() + " " + secondWord.toUpperCase();
        beeline.getActiveSpan()
            .addField("result", result)
            .addField("result-length", result.length());
        return result;
    }

    /**
     * Throws an exception, which demonstrates the capture of error details.
     *
     * @return nothing.
     * @throws RuntimeException to demonstrate error details.
     */
    @RequestMapping("/throw-exception")
    String throwException() {
        throw new RuntimeException("Found an error!");
    }

    /**
     * @return some json data.
     */
    @RequestMapping("/get-captains-data")
    CustomerData getCustomerData() {
        CustomerData customerData = new CustomerData();
        customerData.setName("James T. Kirk");
        customerData.setAge(53);
        return customerData;
    }

    /**
     * Call into an AOP-annotated service, demonstrating declarative child spans.
     *
     * @return a random captain.
     */
    // Since controllers are just Spring beans, they are also subject to AOP, so will get an additional Span here.
    @ChildSpan("PickACaptain")
    @RequestMapping("/pick-a-captain")
    String getARandomCaptain() {
        Span activeSpan = beeline.getActiveSpan().addField("controller-reached", true);

        // Both methods of the "component" are annotated to generate child spans
        String chosenCaptain = component.pickARandomCaptain("Picard", "Kirk", "Pike", "Sisko", "Janeway", "Archer");
        activeSpan.addField("intermediate-result", chosenCaptain);

        String titleAndName = component.prefixWithTitle("Captain", chosenCaptain);
        activeSpan.addField("final-result", chosenCaptain);

        return titleAndName;
    }

    /**
     * The service hands off to other threads and returns early. This shows how async spans can outlive the root.
     */
    @RequestMapping("/run-async-tasks")
    void asyncTask() {
        Span activeSpan = beeline.getActiveSpan().addField("controller-reached", true);
        component.runAsyncTasks();
    }
}
