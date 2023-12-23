include .env
export
# Executables
DOCKER_COMPOSE := docker compose -f docker-compose.yml
DOCKER_COMPOSE_PROD := docker compose -f docker-compose.prod.yml
API_CONTAINER := $(DOCKER_COMPOSE) exec api
API_CONTAINER_PROD := $(DOCKER_COMPOSE_PROD) exec api
WEBAPP_CONTAINER := $(DOCKER_COMPOSE) exec webapp
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
# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #
.PHONY: db/migrations/new
db/migrations/new: ## db/migrations/new name=$1: create a new migration
	@$(eval name ?=)
	@${API_CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: ## db/migrations/up: apply all up migrations
	@${API_CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" up

.PHONY: db/migrations/down
db/migrations/down: ## db/migrations/down: apply all down migrations
	@${API_CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" down

.PHONY: db/migrations/version
db/migrations/version: ## db/migrations/version: print the current migration version
	@${API_CONTAINER} go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./migrations -database="${DATABASE_DSN}" version
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

.PHONY: docker-compose/prod/build
docker-compose/prod/build: ## docker-compose/prod/build: build the production docker images
	@$(DOCKER_COMPOSE_PROD) build --pull --no-cache

.PHONY: docker-compose/prod/up
docker-compose/prod/up: ## docker-compose/prod/up: start the production containers in detached mode
	@$(DOCKER_COMPOSE_PROD) up -d

.PHONY: docker-compose/prod/down
docker-compose/prod/down: ## docker-compose/prod/down: stop the running production containers
	@$(DOCKER_COMPOSE_PROD) down --remove-orphans

# ==================================================================================== #
# Tests
# ==================================================================================== #
.PHONY: test-webapp
test-webapp: ## test: run all tests for webapp
	@${WEBAPP_CONTAINER} go test -v ./...

.PHONY: test-api
test-api: ## test: run all tests for api
	@${API_CONTAINER} go test -v ./...

.PHONY: test
test: test-webapp test-api ## test: run all tests