# Bakery

_Bakery_ is a very simple distributed service I used to learn how to use [Tilt](https://tilt.dev/) and the [OpenTelemetry](https://opentelemetry.io/) Go API.  It consists of 4 parts:  
* 'Reception' - a HTTP server that publishes to a queue  
* 'Baker' - a queue consumer  
* A [RabbitMQ](https://rabbitmq.com/) queue broker  
* A [Redis](https://redis.io/) database used to persist data between the publisher and consumer  

![Basic architecture diagram showing 'reception' using a publishing to a queue and using a database and 'baker' consuming from the queue and using the same database](doc/img/overview.png)

Additionally, when run with Tilt a [Jaeger](https://www.jaegertracing.io/) service is started that 'Baker' and 'Reception' will publish traces to.

## Running  

The following tools and versions were used for developing and testing _Bakery_:

| Tool | Version |
|------|---------|
| Go | v.1.14.2 |
| Tilt | v0.20.7 |  
| Minikube | v1.21.0 |
| Kubernetes | v.1.20.7 |

I recommend using a version manager like [asdf](https://github.com/asdf-vm/asdf) or [goenv](https://github.com/syndbg/goenv) to manage tool versions. The [Tilt documentation](https://docs.tilt.dev/choosing_clusters.html) lists options for local Kubernetes clusters.

Once the required tools are installed, start the service using Tilt:  

```
$ tilt up
Tilt started on http://localhost:10350/
v0.20.7, built 2021-06-10

(space) to open the browser
(s) to stream logs (--stream=true)
(t) to open legacy terminal mode (--legacy=true)
(ctrl-c) to exit
```

Opening the browser should show all the services being started and eventually available:  

![Tilt browser with services started](doc/img/tilt.png)

You can confirm the `reception` service has started by sending a request to the exposed endpoint:

```
$ curl -s http://localhost:8000/ | jq
{
  "message": "reception is attended"
}
```

## Basic Use

Create a new order by sending a `POST` request to the `orders` endpoint:

```
$ curl -s -X POST \
  -d '{"customer": "homer", "pastry": "la bombe"}' \
  http://localhost:8000/orders | jq
{
  "orderKey": "1d5fba984abea5a7de4e2de5d1462bd3"
}
```

Check the status of the order by using the provided order key:

```
$ curl -s http://localhost:8000/orders/1d5fba984abea5a7de4e2de5d1462bd3 | jq
{
  "pastry": "la bombe",
  "customer": "homer",
  "orderTime": "2021-06-19T11:21:24Z",
  "status": "finished"
}
```

Examine traces in Jaeger by browing the exposed endpoint (http://localhost:16686):

![Jaeger UI showing a full collapsed trace](doc/img/jaeger-trace.png)

![Jarger UI showing expanded sections of a trace](doc/img/jaeger-expand.png)

# Clean up

Use tilt to stop the services and delete them from your cluster:  

```
$ tilt down
Beginning Tiltfile execution
Successfully loaded Tiltfile (19.0233ms)
Deleting kubernetes objects:
→ Deployment/baker
→ Deployment/redis
→ Deployment/rabbitmq
→ Deployment/jaeger
→ Deployment/reception
→ Service/jaeger
→ Service/redis
→ Service/rabbitmq
→ ConfigMap/connection-details
```