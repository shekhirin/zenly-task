DOCKER_COMPOSE=docker-compose -p zenly

GRAFANA_URL=http://admin:admin@grafana:3000

.PHONY: protoclean
protoclean:
	@find internal/pb/ -mindepth 1 ! -path '*/mocks/*' -a ! -name "generate.go" -delete

.PHONY: protogen
protogen: protoclean
	@protoc --go_out=./internal/pb \
	--go-grpc_out=./internal/pb \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	-I ./proto \
	$(shell find ./proto -iname "*.proto")

.PHONY: infra-up
infra-up:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml up -d

.PHONY: infra-down
infra-down:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml down

.PHONY: monitoring-up
monitoring-up:
	@$(DOCKER_COMPOSE) -f docker-monitoring.yml up -d

.PHONY: monitoring-down
monitoring-down: grafana-provision
	@$(DOCKER_COMPOSE) -f docker-monitoring.yml down

.PHONY: app-up
app-up:
	@$(DOCKER_COMPOSE) -f docker-compose.yml up --build -d

.PHONY: app-down
app-down:
	@$(DOCKER_COMPOSE) -f docker-compose.yml down

.PHONY: up
up: infra-up monitoring-up app-up

.PHONY: down
down:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml -f docker-monitoring.yml -f docker-compose.yml down

.PHONY: load
load:
	@$(DOCKER_COMPOSE) -f docker-compose.yml build -q
	@docker run --rm --network=zenly_default zenly_zenly ./load -grpc-addr=zenly:8080 -duration=0

.PHONY: mockgen
mockgen:
	@find ~+ -type f -print0 | xargs -0 grep -l "^//go:generate" | sort -u | xargs -L1 -P $$(nproc) go generate

.PHONY: grafana-provision
grafana-provision:
	@docker run --rm --network=zenly_default -v `pwd`/grafana:/home/grafana dwdraju/alpine-curl-jq bash /home/grafana/provision_dashboards.sh $(GRAFANA_URL) /home/grafana/provisioning/dashboards

.PHONY: pumlgen
pumlgen:
	@plantuml diagrams/*.puml
