# zenly-task

[Task description](zenly-task.md)

## Getting Started

### Requirements:
- **Docker**
- **Docker Compose**

Optional:
- **GoLand** for running prepared configurations (load testing tool, application with grpc service)
- **protoc**, **protoc-gen-go-grpc** for generating protobuf implementations
- **plantuml** for generating architecture diagram

### Start
Start all (infrastructure, monitoring, application):
```bash
make up
```
#### OR
Start infrastructure (ZooKeeper, Kafka, NATS):
```bash
make infra-up
```
Start monitoring (Prometheus, Prometheus Exporters, Grafana):
```bash
make monitoring-up
```
Start application:
```bash
make app-up
```

### Stop
Stop all
```bash
make down
```
#### OR
Stop infrastructure:
```bash
make infra-down
```
Stop monitoring:
```bash
make monitoring-down
```
Stop application:
```bash
make app-down
```

## Architecture Overview
![architecture diagram](diagrams/architecture.png "Architecture Diagram")

## Key Points
### Architecture
- NATS uses 1 subject per 1 user
- **gRPC publish stream** [calls enrichers](zenly/enrich.go) and publishes the enriched geolocation to the NATS subject
- Each enricher has 100ms to complete and returns [SetFunc](zenly/enricher/enricher.go) to set the value to
the geolocation if the context isn't timeout-ed yet
- Enrichers are supposed to use [services](zenly/service) to get external data
(e.g. [weather enricher](zenly/enricher/weather.go) uses [weather service](zenly/service/weather/service.go) to get
fake weather data at location)
- **gRPC subscribe stream** subscribes to multiple NATS subjects using [MultiSub](zenly/bus/nats/multisub/multisub.go)
with message handler sending all incoming enriched geolocations to client

### Scalability
All parts of the system can be scaled easily:
- Vertical scaling through spawning more application instances isn't needed, since golang eats all
the resources available
- Horizontal scaling can be done through setting up more pods (in terms of k8s) and a simple load balancer.
[Using a load balancer like Linkerd](https://kubernetes.io/blog/2018/11/07/grpc-load-balancing-on-kubernetes-without-tears/#grpc-load-balancing-on-kubernetes-with-linkerd) 
that allows balancing requests across connections isn't that useful because we use gRPC streams to send and receive
geolocations, thus it's not a good idea to reopen a stream every time
- NATS can be [put into a cluster](https://docs.nats.io/nats-server/configuration/clustering) to have
more throughput and less downtime, then each server in a cluster should be added to the app through `-nats-servers` flag
- More kafka brokers can be added through `-kafka-brokers` flag

### Monitoring:
- Enriching process reports each enricher's time with result (in time / timeout) and
total enriching time with finish reason (complete / timeout)
- Grafana lives at http://localhost:3000/ with username `admin` and password `admin`
