This is an example webhook listener that will hear notifications coming from triggers and print them to STDOUT.

Here's an example of a notification that Honeycomb Triggers would send and this webhook would accept and print out:

```json
{
  "version": "v0.1.0",
  "shared_secret": "would you like to play a game",
  "name": "trig on ttt",
  "status": "TRIGGERED",
  "summary": "Triggered: trig on ttt",
  "description": "Currently greater than threshold value (2) for foo:ble (value 5)",
  "result_url": "",
  "result_groups": [
    {
      "Group": {"foo": "ble"},
      "Result": 5
    },
    {
      "Group": {"foo": "baz"},
      "Result": 1
    },
    {
      "Group": {"foo": "bar"},
      "Result": 1
    }
  ],
  "result_groups_triggered": [
    {
      "Group": {"foo": "ble"},
      "Result": 5
    }
  ]
}
```

The query that's attached to this trigger is:
* breakdown on column `foo`
* alert if `COUNT > 2`
* notify a webhook at `http://myhost.com/notify` with the shared secret `would you like to play a game`

The notification is in the `TRIGGERED` state, which mean it has just crossed the threshold.

The `result_groups` key lists every value for the `foo` column and the counts ofr each one. In this case, `foo` has 3 values: `ble`, `baz`, and `bar`. `bar` and `baz` each only have a `COUNT` of 1, and `ble` has a cound of 5.

The `result_groups_triggered` key only lists the `ble` value becaues it is the only one that is more than 2, the threshold configured in the trigger.
