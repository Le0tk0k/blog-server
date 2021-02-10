ENV_TEST_FILE := .env.test
ENV_TEST = $(shell cat $(ENV_TEST_FILE))

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

.PHONY: up-test-db
up-test-db:
	docker run --rm --env-file=$(ENV_TEST_FILE) -v $(PWD)/build/mysql/my.cnf:/etc/mysql/conf.d/my.cnf  --name blog-server_test_db -d -p 3306:3306 mysql:8.0

.PHONY: down-test-db
down-test-db:
	docker stop blog-server_test_db

.PHONY: test
test:
	$(ENV_TEST) go test -v ./... -count=1

.PHONY: tbls
tbls:
	docker run --rm --net=blog-server_default --env-file=.env.local -v $(PWD):/work k1low/tbls doc -f