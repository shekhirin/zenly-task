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
