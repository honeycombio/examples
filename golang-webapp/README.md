## golang-webapp

Shoutr is an example Golang web application. You can register for accounts, sign
in and shout your opinions on the Internet. It has two tiers: A Golang web app
and a MySQL database.

## Install

Clone the repository into `$GOPATH/src/github.com/honeycombio/examples`.

Create database in MySQL.

```
$ mysql -uroot -e 'create database shoutr;'
```

Run app.

```
$ export HONEYCOMB_WRITEKEY=<writekey>
$ go run main.go
```

## Run in Docker

The whole webapp can be run in Docker (Compose).

Set your [Honeycomb write key](https://ui.honeycomb.io/account) to
`HONEYCOMB_WRITEKEY`, or edit the `docker-compose.yml`. The `shoutr` database in
MySQL will be created automatically.

```
$ export HONEYCOMB_WRITEKEY=<writekey>
```

Then:

```
$ docker-compose up
```

## Event Fields

| **Name** | **Description** | **Example Value** |
| --- | --- | --- |
| `flash.value` | Contents of the rendered flash message | `Your shout is too long!` |
| `request.content_length`| Length of the content (in bytes) of the sent HTTP request | `952` |
| `request.host` | Host the request was sent to | `localhost` |
| `request.method` | HTTP method | `POST` |
| `request.path` | Path of the request | `/shout` |
| `request.proto` | HTTP protocol version | `HTTP/1.1` |
| `request.remote_addr` | The IP and port that answered the request  | `172.18.0.1:40484` |
| `request.user_agent`| User agent | `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36` |
| `response.status_code` | Status code written back to user | `200` |
| `runtime.memory_inuse` | Amount of memory in use (in bytes) by the whole process | `4,971,776` |
| `runtime.num_goroutines` | Number of goroutines in the process | `7` |
| `shout.content` | Content of the user's comment | `Hello world!` |
| `shout.content_length` | Length (in characters) of the user's comment | `80` |
| `system.hostname` | System hostname | `1ba87a98788c` |
| `timers.total_time_ms` | The total amount of time the request took to serve | `180` |
| `timers.mysql_insert_user_ms` | The time the `INSERT INTO users` query took | `50` |
| `user.id`| User ID | `2` |

## Example Queries
