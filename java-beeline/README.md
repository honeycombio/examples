# beeline-java example web app

This application demonstrates a simple Java Beeline instrumentation of 
a [Spring Boot](https://projects.spring.io/spring-boot/) application, including adding child spans to traces and custom fields to spans.

## Run locally

### Required configuration

Set your HoneyComb write key and the dataset you want to report to by using the 
[application.properties](src/main/resources/application.properties) file. 

```
# (Required) Dataset to send the Events/Spans to
honeycomb.beeline.dataset                :testDataset

# (Required) Your honeycomb account write key
honeycomb.beeline.write-key              :testKey
```

This application requires Java 8.

### Run commands

```sh
mvn install
mvn spring-boot:run
```