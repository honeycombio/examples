## python-api

The Python API example shows how to do Honeycomb instrumentation with a Flask app.

Events are created per-HTTP-request (following the Honeycomb _one event per unit of work_
model) using the `@app.before_request` decorator and sent after the request using the
`@app.after_request` decorator. Because the events are stored on `g`, Flask's thread 
local storage, the events contain a variety of default properties describing the request
as well as custom fields that can be added by any handler using `g.ev.add_field`.

## Run Locally

Set your [Honeycomb write key](https://ui.honeycomb.io/account) to
`HONEYCOMB_WRITEKEY`.

MySQL and Flask must be installed (`pip install flask`).

First, create the database used by the app.

```
$ mysql -uroot -e 'create database example-python-api;'
```

Then run the Python app on `localhost`:

```
$ FLASK_APP=app.py flask run
```

A basic REST API for todos is exposed on port 5000.

```
$ curl \
    -H 'Content-Type: application/json' \
    -X POST -d '{"description": "Walk the dog", "due": 1518816723}' \
    localhost:5000/todos/
...

$ curl localhost:5000/todos/
[
  {
    "completed": false,
    "description": "Walk the dog",
    "due": "Fri, 16 Feb 2018 21:32:03 GMT",
    "id": 1
  }
]

$ curl -X PUT \
    -H 'Content-Type: application/json' \
    -d '{"description": "Walk the cat"}' \
    localhost:5000/todos/1/
{
  "completed": false,
  "description": "Walk the cat",
  "due": "Fri, 16 Feb 2018 21:32:03 GMT",
  "id": 1
}

$ curl -X DELETE localhost:5000/todos/1/
{
  "id": 1,
  "success": true
}

$ curl localhost:5000/todos/
[]
```

## Run in Docker

This example can be run in Docker (Compose).

```
$ docker-compose build && docker-compose up -d
```

## Event Fields

| **Name** | **Description** | **Example Value** |
| --- | --- | --- |
| `errors.message` | Message in the error encountered, if applicable | `undefined` |
| `request.endpoint` | Endpoint requested | `/todos/` |
| `request.method` | HTTP method | `POST` |
| `request.path` | Request path | `/todos/` |
| `request.python_function` | Python function serving the request | `index` |
| `request.url_pattern` |` Underlying routing pattern of the URL | `/todos/<int:todo_id>/` |
| `request.user_agent` | User agent for the request | `curl/7.54.0` |
| `request.user_agent.browser` | Web browser the request was served to | `chrome` |
| `request.user_agent.platform` | OS of the user agent | `macos` |
| `request.user_agent.string` | Literal user agent string | `curl/7.54.0` |
| `request.user_agent.version` | Version of the user agent | `64.0.3282.186` |
| `response.status_code` | HTTP status code of the response | 404 |
| `timers.db.delete_todo` | Time in milliseconds for DB call to delete a todo | 23 |
| `timers.db.insert_todo_ms` | Time in milliseconds for DB call to insert a todo | 50 |
| `timers.db.select_all_todos` | Time in millisconds for DB call to select all todos | 11 |
| `timers.db.select_todo` | Time in milliseconds for DB call to select a todo | 4 |
| `timers.db.update_todo` | Time in milliseconds for DB call to update a todo | 50 |
| `timers.flask_time_ms` | Total time in milliseconds Flask spent serving the request | 75 |
| `todo.id` | ID of the associated TODO | 1 |

## Example Queries

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/python-api-q1.png)
