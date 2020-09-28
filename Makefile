DOCKER_COMPOSE=docker-compose -p zenly

.PHOY: proto-clean
proto-clean:
	@find zenly/pb/ -mindepth 1 -delete

.PHONY: proto-gen
proto-gen: proto-clean
	@protoc --go_out=./zenly/pb \
	--go-grpc_out=./zenly/pb \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	-I ./proto \
	proto/*.proto proto/**/*.proto

.PHONY: infra-up
infra-up:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml up -d

.PHONY: infra-up
infra-down:
	@$(DOCKER_COMPOSE) -f docker-infrastructure.yml down

.PHONY: up
up:
	@$(DOCKER_COMPOSE) -f docker-compose.yml up --build -d

.PHONY: down
down:
	@$(DOCKER_COMPOSE) -f docker-compose.yml down
