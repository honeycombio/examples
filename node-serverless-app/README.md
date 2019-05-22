# node-serverless-app

This example illustrates a Lambda function instrumented with Honeycomb's [Beeline for Node](https://docs.honeycomb.io/getting-data-in/javascript/beeline-nodejs/).

It contains examples of:

- Baseline Beeline usage (not in an Express app)
- Capture of custom metadata on Beeline-generated events
- Definition of custom spans to augment traces
- Continuing a propagated trace (if applicable), e.g. if the client initiated a trace via another Honeycomb Beeline
- Error handling with a trace

## Usage:

Find your Honeycomb API key at https://ui.honeycomb.io/account, then make it available to your Lambda function.

After deploying your Lambda function, you may curl it with a payload:

```
curl -X POST $LAMBDA_URL -d 'foo bar baz payload'
```

## Generated traces:

The example code should produce the below `lambda` spans in a given trace:

![lambda spans](/images/trace.png)
