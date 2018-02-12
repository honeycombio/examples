## honeytail-nginx

Starting with Honeycomb instrumentation "at the edge" (i.e., with your reverse
proxies or load balancers) can allow you to quickly get valuable data into
Honeycomb. You can begin gaining visibility into your systems quickly this way,
and gradually work your way inward (consequently gaining more power) by adding
native code instrumentation later.

This example demonstrates this concept by using nginx as a reverse proxy to the
[Python API
example](https://github.com/honeycombio/examples/tree/master/python-api), and
ingesting the nginx access logs as Honeycomb events.

## Run Natively

Run the Python API example linked above.

Use the provided `nginx.conf` for your local nginx config. You may need to
update the `proxy_pass` to pass to `localhost` instead of `api`. Then:

```
$ honeytail --debug \
    --parser=nginx \
    --dataset=examples.honeytail-nginx \
    --writekey=$HONEYCOMB_WRITEKEY \
    --nginx.conf=/etc/nginx/nginx.conf \
    --nginx.format=honeytail \
    --file=/var/log/honeytail/access.log
```

## Run in Docker

```
$ docker-compose build && docker-compose up -d
$ curl localhost/api
{"up": true}
$ curl \
    -H 'Content-Type: application/json' \
    -X POST -d '{"description": "Walk the dog", "due": 1518816723}' \
    localhost/api/todos/
$ curl localhost/api/todos/
[
  {
    "completed": false,
    "description": "Walk the dog",
    "due": "Fri, 16 Feb 2018 21:32:03 GMT",
    "id": 1
  }
]
... etc ...
```

## Event Fields

| **Name** | **Description** | **Example Value** |
| --- | --- | --- |
| `body_bytes_sent` | # of bytes in the HTTP response body sent to the client | 157 |
| `bytes_sent` | # of bytes sent to the client | 24  |
| `host` | Hostname of the nginx server responding to the request | `endophage` |
| `http_referer` | HTTP referer header | `https://developer.mozilla.org/en-US/docs/Web/JavaScript` |
| `http_user_agent` | User agent of the request | `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36` |
| `remote_addr` | Client address | `172.19.0.1` |
| `request` | Full original request line | `GET /favicon.ico HTTP/1.1` |
| `request_length` | Size of request in bytes | 456 |
| `request_method` | HTTP request method | `POST` |
| `request_path` | URL of the request | `/api/todos/` |
| `request_pathshape` | "Shape" of the request path | `/api/todos/` |
| `request_protocol_version` | HTTP version | `HTTP/1.1` |
| `request_shape` | "Shape" of the request | `/api/todos/` |
| `request_time` | Amount of time it took to serve the request in seconds | 250 |
| `request_uri` | URI for the request | `/api/todos/`` |
| `server_name` | Name of the server serving the request | `localhost` |
| `status` | HTTP status code returned | 404 |

## Example Queries

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/honeytail-nginx-q1.png)
