include .env
export
# Executables
DOCKER_COMPOSE := docker compose -f docker-compose.yml
CONTAINER := $(DOCKER_COMPOSE) exec api

# Database
DATABASE_DSN := mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST))/$(MYSQL_DATABASE)?charset=utf8mb4

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
.DEFAULT_GOAL = help
## help: print this help message
.PHONY: help
help:
	@grep -E '(^[a-zA-Z0-9\./_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## start: cold-start the application
.PHONY: start
start:
	make docker-compose/up
	make db/migrations/up

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #
.PHONY: db/migrations/new
db/migrations/new: ## db/migrations/new name=$1: create a new migration
	@$(eval name ?=)
	@${CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: ## db/migrations/up: apply all up migrations
	@${CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" up

.PHONY: db/migrations/down
db/migrations/down: ## db/migrations/down: apply all down migrations
	@${CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" down

.PHONY: db/migrations/version
db/migrations/version: ## db/migrations/version: print the current migration version
	@${CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" version
# ==================================================================================== #
# DOCKER_COMPOSE
# ==================================================================================== #
.PHONY: docker-compose/build
docker-compose/build: ## docker-compose/build: build the docker images
	@$(DOCKER_COMPOSE) build --pull --no-cache

.PHONY: docker-compose/up
docker-compose/up: ## docker-compose/up: start the containers in detached mode
	@$(DOCKER_COMPOSE) up -d

.PHONY: docker-compose/down
docker-compose/down: ## docker-compose/down: stop the running containers
	@$(DOCKER_COMPOSE) down --remove-orphans

# ==================================================================================== #
# Tests
# ==================================================================================== #
.PHONY: tests
tests: ## test: run all tests
	@${CONTAINER} go test -v ./...
