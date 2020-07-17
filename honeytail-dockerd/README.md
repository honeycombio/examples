## honeytail-dockerd

The Docker daemon logs using [logrus](https://brandur.org/logfmt), and emits
key-value pair log lines similar to the [logfmt](https://brandur.org/logfmt)
format. e.g.:

```
time="2018-04-10T23:41:04.155700308Z" level=debug msg="EnableService 49f92179af6b93864dd69de743312e196efc98be131fc35f19aea6f1dfa2e3af DONE"
time="2018-04-10T23:41:04.159635917Z" level=debug msg="bundle dir created" bundle=/var/run/docker/containerd/49f92179af6b93864dd69de743312e196efc98be131fc35f19aea6f1dfa2e3af module=libcontainerd namespace=moby root=/var/lib/docker/overlay2/e206dbff938f529daf7e5f385d1bdcc945b688d7c13bd5f073edb861449969b0/merged
time="2018-04-10T23:41:04Z" level=debug msg="event published" module="containerd/containers" ns=moby topic="/containers/create" type=containerd.events.ContainerCreate
time="2018-04-10T23:41:04Z" level=info msg="shim docker-containerd-shim started" address="/containerd-shim/moby/49f92179af6b93864dd69de743312e196efc98be131fc35f19aea6f1dfa2e3af/shim.sock" debug=true module="containerd/tasks" pid=1676
time="2018-04-10T23:41:04Z" level=debug msg="registering ttrpc server"
time="2018-04-10T23:41:04Z" level=debug msg="serving api on unix socket" socket="[inherited from parent]"
time="2018-04-10T23:41:04.281146967Z" level=debug msg="sandbox set key processing took 49.81037ms for container 49f92179af6b93864dd69de743312e196efc98be131fc35f19aea6f1dfa2e3af"
```

These structured logs can be parsed into Honeycomb events using
[Honeytail](https://honeycomb.io/docs/connect/agent/).

This should work on Docker for Mac and Windows. To make the example run on
Linux, the location of the Docker log file will need to be updated in
`docker-compose.yml`.

**Important:** While arbitrary log lines can be transformed into Honeycomb
events as demonstrated here, following the "one-event-per-unit-of-work" model
will generally have better results for observability of your natively
instrumented apps.

## Run Natively

Set your [Honeycomb write key](https://ui.honeycomb.io/account) to
`HONEYCOMB_WRITEKEY`.

```
$ export HONEYCOMB_WRITEKEY=<writekey>
```

Then invoke `honeytail`:

```
$ honeytail --debug \
    --parser=keyval \
    --dataset=examples.honeytail-dockerd \
    --writekey=$HONEYCOMB_WRITEKEY \
    --keyval.timefield=time \
    --file=/var/log/docker.log
```

You may need to update the location of the Docker log file.

## Run in Docker

A `docker-compose.yml` is provided to run this example.

```
$ docker-compose build && docker-compose up
```

## Event Fields

| **Name** | **Description** | **Example Value** |
| --- | --- | --- |
| `address` | The socket `dockerd` is connecting to | `/var/run/docker/containerd/docker-containerd.sock` |
| `bundle` | The "bundle" `dockerd` is working with | `/var/run/docker/containerd/5cd9224ae7761c052173f34bf40226cf24a472d956bd8c3d2d2f2446cbbf375a` |
| `debug` | Whether or not this is a debug mode message | `true` |
| `error_type` | The type of error encountered (if present) | `*errors.fundamental` |
| `id` | The container or image ID associated with this log line | `30660078ad577d3eebfbcbb8e1a073e214180fb23627bd2f051d0830af44fd4a` |
| `level` | The log level (warning, info, debug, etc.) | `warning` |
| `module` | The subsection of `dockerd` this message is from | `containerd/tasks` |
| `msg` | The log message itself | `Calling GET /_ping` |
| `namespace` | The namespace with which the log lines are associated | `plugins.moby` |
| `ns` | The namespace with which the log lines are associated | `moby` |
| `pid` | The process ID of an exec-ed process, e.g., the containerd shim | `9055` |
| `root` | The root of the union filesystem used for the container | `/var/lib/docker/overlay2/9f65055258cea318a5aaeac58b6dbf1b3fde7ad066cb2711e7046b392d8903bd/merged` |
| `socket` | Which socket is in usage during this log line  | `[inherited from parent]` |
| `topic` | The action under operation | `/tasks/start` |
| `type` | Type of action being taken | `io.containerd.snapshotter.v1` |

## Example Queries

![](https://raw.githubusercontent.com/honeycombio/examples/main/_internal/honeytail-dockerd-q1.png)
