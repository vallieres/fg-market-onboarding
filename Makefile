DEFAULT_GOAL: help


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
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	go install gotest.tools/gotestsum@latest
	go install github.com/air-verse/air@latest

.PHONY: start
start: ## Runs the service
	air

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
