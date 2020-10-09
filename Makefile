DOCKER_COMPOSE=docker-compose -p zenly

GRAFANA_URL=http://admin:admin@grafana:3000

.PHOY: protoclean
protoclean:
	@find zenly/pb/ -mindepth 1 -delete

.PHONY: protogen
protogen: protoclean
	@protoc --go_out=./zenly/pb \
	--go-grpc_out=./zenly/pb \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	-I ./proto \
	$(shell find ./proto -iname "*.proto")

.PHONY: infra-up
infra-up:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml up -d

.PHONY: infra-up
infra-down:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml down

.PHONY: monitoring-up
monitoring-up:
	@$(DOCKER_COMPOSE) -f docker-monitoring.yml up -d

.PHONY: monitoring-up
monitoring-down: grafana-provision
	@$(DOCKER_COMPOSE) -f docker-monitoring.yml down

.PHONY: up
up:
	@$(DOCKER_COMPOSE) -f docker-compose.yml up --build -d

.PHONY: down
down:
	@$(DOCKER_COMPOSE) -f docker-compose.yml down

.PHONY: mockgen
mockgen:
	@find ~+ -type f -print0 | xargs -0 grep -l "^//go:generate" | sort -u | xargs -L1 -P $$(nproc) go generate

.PHONY: grafana-provision
grafana-provision:
	@docker run --rm --network=zenly_default -v `pwd`/grafana:/home/grafana dwdraju/alpine-curl-jq bash /home/grafana/provision_dashboards.sh $(GRAFANA_URL) /home/grafana/provisioning/dashboards
