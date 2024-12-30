DEFAULT_GOAL: help

DBCONNECT := "fgonboard:fgonboard@/fgonboard?parseTime=true"

ifneq (,$(wildcard ./.env))
	include .env
	export
endif


.PHONY: help
help: ## Print this message and exit.
	@awk 'BEGIN {FS = ":.*?## "} /^[0-9a-zA-Z_-]+:.*?## / {printf "\033[36m%s\033[0m : %s\n", $$1, $$2}' $(MAKEFILE_LIST) \
		| sort \
		| column -s ':' -t

install: ## Install required software
	brew install agrinman/tap/tunnelto
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	go install gotest.tools/gotestsum@latest
	go install github.com/air-verse/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: up
up: ## Starts the database
	docker-compose up --detach
	sleep 7
	goose -dir ./app/db mysql $(DBCONNECT) up

.PHONY: down
down: ## Stops the database
	docker-compose stop

.PHONY: db-up
db-up: ## Run migrations up to last version
	goose -dir ./app/db mysql $(DBCONNECT) up

.PHONY: db-status
db-status: ## Shows which migration version we are at
	goose -dir ./app/db mysql $(DBCONNECT) status

.PHONY: db-down
db-down: ## Run migrations down
	goose -dir ./app/db mysql $(DBCONNECT) down

.PHONY: db-add
db-add: ## Add a new migration
	@cd app/db; read -p "What is the new migration about (slug)?: " migname; \
	goose create $$migname sql
	@cd ..

.PHONY: start
start: ## Runs the service
	docker-compose up --detach
	sleep 15
	goose -dir ./app/db mysql $(DBCONNECT) up
#	export $(grep -v -e '^$' .env | grep -v -e '^#' | xargs -0)
	FGONBOARDING_ENVIRONMENT=local FGONBOARDING_SERVER_PORT=443 air

.PHONY: start-tunnel
start-tunnel: ## Runs the service and opens tunnel
	docker-compose up --detach
	sleep 15
	goose -dir ./app/db mysql $(DBCONNECT) up
#	export $(grep -v -e '^$' .env | grep -v -e '^#' | xargs -0)
	FGONBOARDING_ENVIRONMENT=dev FGONBOARDING_SERVER_PORT=8004 air & ssh -R 80:localhost:8004 serveo.net && fg

.PHONY: format
format: ## Format all the files
	goimports -local github.com/vallieres/fg-market-onboarding -w .
	gofumpt -l -w .
	go mod tidy

.PHONY: lint
lint: ## Lint the Go code
	golangci-lint run -c .golangci.yml

.PHONY: test
test: ## Test
	@echo $(HYPR_DB_SERVER_URL)
	go run main.go

define slugify
    echo "$1" | iconv -t ascii//TRANSLIT | sed -r s/[~\^]+//g | sed -r s/[^a-zA-Z0-9]+/-/g | sed -r s/^-+\|-+$//g | tr A-Z a-z
endef
