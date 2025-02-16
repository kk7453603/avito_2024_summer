# Импортирует переменные окружения из .env для доступа к ним из Makefile
ifneq (,$(wildcard .env))
    include .env
    export
endif

DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: m-up
# migrate up
m-up:
	@docker-compose run --rm migrate -path ./migrations -database "$(DSN)" up

.PHONY: m-down
# migrate down
m-down:
	@docker-compose run --rm migrate -path ./migrations -database "$(DSN)" down -all

.PHONY: m-status
# migrate version (status)
m-status:
	@docker-compose run --rm migrate -path ./migrations -database "$(DSN)" version

.PHONY: d-up
# service improvement
d-up:
	@docker-compose up -d

.PHONY: d-up-b
# service improvement with rebuild
d-up-b:
	@docker-compose up -d --build

.PHONY: d-down
# service closure
d-down:
	@docker-compose down

.PHONY: d-down-v
# service closure with data deletion
d-down-v:
	@docker-compose down -v

.PHONY: d-up-app
# raising the service without auto-migration, if data remains
d-up-app:
	@docker-compose up -d postgres app

.PHONY: lint
# linter start
lint:
	@golangci-lint run

COVER_PKG_LIST ?= ./internal/db ./internal/hasher ./internal/modules/authentication ./internal/modules/jwt_token_manager \
./internal/modules/buy_item ./internal/modules/transaction ./internal/modules/user_info \
./internal/server ./internal/server/handlers

.PHONY: tests
# running all tests except integration tests
tests:
	@go test -v -count=1 $(COVER_PKG_LIST)

.PHONY: tests-integration
# raises the test database, waits for initialization and performs integration and other tests, then deletes the test container
tests-integration:
	@echo "-> Starting PostgreSQL and applying migrations..."
	@docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml up -d
	@echo "-> Running integration tests..."
	@go test -tags=integration -v -count=1 $(COVER_PKG_LIST) \
     	|| (docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml down -v && exit 1)
	@docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml down -v

.PHONY: cover
# viewing the result of code coverage by tests in html form
cover:
	@go test -short -count=1 -coverprofile=coverage.out $(COVER_PKG_LIST)
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out

.PHONY: cover-integration
# view code coverage results in html format including integration tests
TEST_CONTAINER_NAME = avitobta2025-merchsop-tests
cover-integration:
	@echo "-> Starting PostgreSQL and applying migrations..."
	@docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml up -d
	@echo "-> Running integration tests..."
	@go test -tags=integration -short -count=1 -coverprofile=coverage.out $(COVER_PKG_LIST) \
 		|| (docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml down -v && exit 1)
	@docker-compose -p $(TEST_CONTAINER_NAME) -f docker-compose-test.yml down -v
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
