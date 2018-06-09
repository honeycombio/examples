package io.honeycomb.libhoney.example.webapp.instrumentation;

import io.honeycomb.libhoney.Event;
import io.honeycomb.libhoney.HoneyClient;
import nl.basjes.parse.useragent.UserAgent;
import nl.basjes.parse.useragent.UserAgentAnalyzer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.handler.HandlerInterceptorAdapter;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.util.Map;

@Component
public class HoneycombHandlerInterceptor extends HandlerInterceptorAdapter {
    private static final String START_TIME_KEY = "start-time";
    private static final UserAgentAnalyzer USER_AGENT_ANALYZER = UserAgentAnalyzer.newBuilder()
                                                                                    .withField(UserAgent.AGENT_NAME)
                                                                                    .withField(UserAgent.OPERATING_SYSTEM_NAME)
                                                                                    .withField(UserAgent.AGENT_VERSION)
                                                                                    .build();

    @Autowired
    private HoneycombContext honeycombContext;
    @Autowired
    private HoneyClient honeyClient;

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        final Map<String, Object> contextualData = honeycombContext.getContextualData();
        final long startTime = System.currentTimeMillis();
        contextualData.put(START_TIME_KEY, startTime);

        final Event requestEvent = honeycombContext.getEvent();
        final String userAgentHeaderValue = request.getHeader(HttpHeaders.USER_AGENT);
        requestEvent.addField("request.user_agent.string", userAgentHeaderValue);
        if (userAgentHeaderValue != null) {
            final UserAgent userAgent = USER_AGENT_ANALYZER.parse(userAgentHeaderValue);
            requestEvent.addField("request.user_agent.browser", userAgent.get(UserAgent.AGENT_NAME).getValue());
            requestEvent.addField("request.user_agent.platform", userAgent.get(UserAgent.OPERATING_SYSTEM_NAME).getValue());
            requestEvent.addField("request.user_agent.version", userAgent.get(UserAgent.AGENT_VERSION).getValue());
        }
        requestEvent.addField("request.path", request.getRequestURI());
        requestEvent.addField("request.method", request.getMethod());
        return true;
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) {
        final Long startTime = (Long) honeycombContext.getContextualData().remove(START_TIME_KEY);
        if (startTime != null) {
            final Event event = honeycombContext.getEvent();
            final long endTime = System.currentTimeMillis();
            event.addField("timers.total_time_ms", endTime - startTime);
            event.addField("response.status_code", response.getStatus());
            event.send();
        }
    }
}
