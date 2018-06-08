package io.honeycomb.libhoney.example.webapp.configuration;

import org.springframework.boot.context.properties.ConfigurationProperties;

import java.util.HashMap;
import java.util.Map;

@ConfigurationProperties("honeycomb")
public class HoneycombProperties {
    private String apiHost;
    private int sampleRate;
    private String dataset;
    private String writeKey;
    private Map<String, String> globalFields = new HashMap<>();

    public String getApiHost() {
        return apiHost;
    }

    public void setApiHost(String apiHost) {
        this.apiHost = apiHost;
    }

    public int getSampleRate() {
        return sampleRate;
    }

    public void setSampleRate(int sampleRate) {
        this.sampleRate = sampleRate;
    }

    public String getDataset() {
        return dataset;
    }

    public void setDataset(String dataset) {
        this.dataset = dataset;
    }

    public String getWriteKey() {
        return writeKey;
    }

    public void setWriteKey(String writeKey) {
        this.writeKey = writeKey;
    }

    public Map<String, String> getGlobalFields() {
        return globalFields;
    }

    public void setGlobalFields(Map<String, String> globalFields) {
        this.globalFields = globalFields;
    }
}
