.PHOY: proto-clean
proto-clean:
	@rm -r internal/pb/*

.PHONY: proto-gen
proto-gen: proto-clean
	@protoc --go_out=./internal/pb \
	--go-grpc_out=./internal/pb \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	-I ./proto \
	proto/*.proto proto/**/*.proto

.PHONY: infra-up
infra-up:
	@docker-compose -f docker-infrastructure.yml up -d

.PHONY: infra-up
infra-down:
	@docker-compose -f docker-infrastructure.yml down

.PHONY: up
up:
	@docker-compose -f docker-compose.yml up --build -d

.PHONY: down
down:
	@docker-compose -f docker-compose.yml down
