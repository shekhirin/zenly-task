DOCKER_COMPOSE=docker-compose -p zenly

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

.PHONY: up
up:
	@$(DOCKER_COMPOSE) -f docker-compose.yml up --build -d

.PHONY: down
down:
	@$(DOCKER_COMPOSE) -f docker-compose.yml down

.PHONY: mockgen
mockgen:
	@find ~+ -type f -print0 | xargs -0 grep -l "^//go:generate" | sort -u | xargs -L1 -P $$(nproc) go generate
