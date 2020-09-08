Honeycomb triggers can specify a webhook as the notification target. When configured in this way, Honeycomb will send an HTTP POST to the URL specified in the trigger configuration. That POST will include:

- A shared secret token for authentication in the `X-Honeycomb-Webhook-Token` header
- The results of the trigger as JSON in the body.

This is an example webhook listener that will hear notifications coming from triggers and print them to STDOUT. It is instrumented with the Honeycomb beeline so you can see what your webhook is doing!

Here's an example of a notification that Honeycomb Triggers would send and this webhook would accept and print out:

```json
{
  "version": "v0.1.0",
  "id": "abdcefg",
  "name": "trig on ttt",
  "trigger_description": "To troubleshoot, please look up the steps in our runbook",
  "trigger_url": "https://ui.honeycomb.io/team/datasets/dataset/triggers/abdcefg",
  "status": "TRIGGERED",
  "summary": "Triggered: trig on ttt",
  "description": "Currently greater than threshold value (2) for foo:fooOOOddd (value 5)",
  "operator": "greater than",
  "threshold": 2,
  "result_url": "",
  "result_groups": [
    {
      "Group": { "foo": "fooOOOddd" },
      "Result": 5
    },
    {
      "Group": { "foo": "hungry" },
      "Result": 1
    },
    {
      "Group": { "foo": "chompy" },
      "Result": 1
    }
  ],
  "result_groups_triggered": [
    {
      "Group": { "foo": "fooOOOddd" },
      "Result": 5
    }
  ]
}
```

Note:

- The field `id` refers to the trigger id that is used to specify this specific trigger in the API.
- The field `trigger_url` is the hyperlink to that trigger in the UI.

The query that's attached to this trigger is:

- breakdown on column `foo`
- alert if `COUNT > 2`
- notify a webhook at `http://myhost.com:8090/notify` with the shared secret `would you like to play a game`

The notification is in the `TRIGGERED` state, which mean it has just crossed the threshold.

The `result_groups` key lists every value for the `foo` column and the counts ofr each one. In this case, `foo` has 3 values: `fooOOOddd`, `hungry`, and `chompy`. `chompy` and `hungry` each only have a `COUNT` of 1, and `fooOOOddd` has a cound of 5.

The `result_groups_triggered` key only lists the `fooOOOddd` value becaues it is the only one that is more than 2, the threshold configured in the trigger.

## Install

Clone the repository into \$GOPATH/src/github.com/honeycombio/examples.

Run app.

    $ go run main.go

Now configure Honeycomb to send a notification to the running listener! Within the Integrations section of your Honeycomb team settings, create a Trigger Recipient with the parameters: Provider: `Webhook`, Webhook URL: `http://example.com:8090/notify`. Create a trigger and add this Webhook as a Recipient. When your trigger fires, Honeycomb will notify this Webhook! (Tools like [ngrok](https://ngrok.com/) can help your local process be available via the public web!)
