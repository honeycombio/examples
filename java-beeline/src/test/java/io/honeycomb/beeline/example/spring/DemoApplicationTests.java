package io.honeycomb.beeline.example.spring;

import io.honeycomb.beeline.tracing.ids.TraceIdProvider;
import io.honeycomb.beeline.tracing.ids.UUIDTraceIdProvider;
import io.honeycomb.beeline.tracing.propagation.HttpHeaderV1PropagationCodec;
import io.honeycomb.beeline.tracing.propagation.Propagation;
import io.honeycomb.beeline.tracing.propagation.PropagationContext;
import io.restassured.RestAssured;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.web.server.LocalServerPort;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.junit4.SpringRunner;

import java.util.Collections;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.startsWith;

@RunWith(SpringRunner.class)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
public class DemoApplicationTests {

    @LocalServerPort
    int port;

    @Before
    public void setUp() {
        RestAssured.port = port;
    }

    @Test
    public void checkThatContextInitializesCorrectly() {

    }

    @Test
    public void checkThatRequestWorkd() {
        TraceIdProvider idProvider = UUIDTraceIdProvider.getInstance();
        final String traceId = idProvider.generateId();
        final String spanId = idProvider.generateId();
        final PropagationContext context = new PropagationContext(traceId, spanId, null, Collections.singletonMap("env", "test"));
        final String encode = Propagation.honeycombHeaderV1().encode(context).get();

        given()
            .header(HttpHeaderV1PropagationCodec.HONEYCOMB_TRACE_HEADER, encode)

            .when().get("/forward-to/{forwardingPath}", "pick-a-captain")

            .then().body(startsWith("Captain"));
    }

}

