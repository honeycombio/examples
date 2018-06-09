## Kubernetes Envoy Tracing

This shows a demonstration of using Envoy's Zipkin instrumentation to get
distributed tracing in Honeycomb with a minimal amount of changes app code.

It is a very similar variant of the [Envoy Zipkin
example](https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/zipkin_tracing)
ported to run on Kubernetes. Instead of sending to Zipkin, Envoy is configured
to send the OpenTracing compatible data to the [Honeycomb OpenTracing
proxy](https://github.com/honeycombio/honeycomb-opentracing-proxy).

## Running the Example

Ensure that the Honeycomb write key is set (you can check with `kubectl get
secrets`):

```
$ kubectl create secret generic honeycomb-writekey --from-literal=key=<WRITEKEY>
```

Then, from this directory:

```
$ kubectl apply -f demo.yaml
```

All the needed images, etc. should be downloaded automatically.

If you need to start over completely from scratch:

```
$ kubectl delete -f demo.yaml
$ kubectl apply -f demo.yaml
```

## Using the app

This diagram from the Envoy documentation gives a good general idea of the
structure of the app. There is a front Envoy which receives all inbound
requests, and forwards them all to `service1`, which then in turn calls
`service2` multiple times. All service-to-service communication happens via
Envoys.

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy-example-arch.svg)

- Envoy is listening for ingress on port 80 in each pod
- For egress, `service1` hits `localhost:9000` (Envoy will listen on this port
  in the pod and proxy outbound requests to the Envoy for `service2`)
- Any service-to-service hop will forward the headers used for tracing (you can
  see this in the Flask app source code) - this is the only app code change
  needed for tracing instrumentation, Envoy will handle the rest.

Once the pods are running, this command will forward port 80 from the front
Envoy pod to `localhost:8000`:

```
$ kubectl port-forward $(kubectl get --no-headers=true pods -o name | grep front)
8000:80
```

You can get the service to echo back a string (wrapped in some HTML) by using
`/echo/<string>`:

```
$ curl localhost:8000/echo/friend
<img src="https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy.svg" height="100" />
<pre><code>Hello friend!

Served by:
service 1
pod: service1-68d8d4ff9f-vv2vp</code></pre>
```

or in browser:

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy-reply.png)

If you query Honeycomb you should be able to see your exact request given the
right parameters. `BREAK DOWN` by `serviceName` (i.e., egress pod) and
`http.url` is a fun one.

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy-heatmap.png)

Clicking on the graph as suggested will render a full trace view. `serviceName`
is the ID of the Pod used for egress for each step of the trace. `name` is the
service (or address, in some rare cases) it called. As you can see, `service1`
calls `service2` five times, and `/stage/3` is usually the slowest.

![](https://raw.githubusercontent.com/honeycombio/examples/master/_internal/envoy-traceview.png)

## Building the Docker image

The Envoy containers and Honeycomb OpenTracing proxy use pre-packaged upstream
images from Docker Hub.

If you need to build the Docker image for the Flask app, DON'T PANIC! It's
fairly simple.

```
$ docker build -t honeycombio/dockercon-2018-flaskapp:<newtag> -f service.Dockerfile .
...
$ docker push honeycombio/dockercon-2018-flaskapp:<newtag>
...
```

You will need to update the Kubernetes YAML (`image:`) to use the image with the
new tag and re-`apply` to get Kube to pick up your changes.

## Generating data

If you have `curl` installed, you can use the script `generate-load.sh` in this
directory to generate some fake load:

```
$ ./generate-load.sh
...
```

With the caveat that this might have implications for the scenario demonstrated
here (i.e., `service2/stage/3` is the slow part of the trace). Still, it'll
likely result in some damn pretty `HEATMAP`s.
