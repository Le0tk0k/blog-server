.PHONY: up
up:
	docker-compose -f docker-compose.local.yml up -d --build

.PHONY: down
down:
	docker-compose -f docker-compose.local.yml down

.PHONY: logs
logs:
	docker-compose -f docker-compose.local.yml logs

.PHONY: ps
ps:
	docker-compose -f docker-compose.local.yml ps

.PHONY: lint
lint:
	golangci-lint run --enable=golint,gosec,prealloc,gocognit,goimports