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
