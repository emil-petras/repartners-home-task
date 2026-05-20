.PHONY: help build run test clean docker-up docker-down fmt deps

help: ## show available targets
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## build the application binary
	@go build -o bin/server ./cmd/server

run: ## run the application locally
	@go run ./cmd/server/main.go

test: ## run all tests
	@go test -v ./...

clean: ## clean build artifacts and test databases
	@rm -rf bin/
	@rm -f *.db test_*.db

docker-up: ## start with docker compose
	@docker-compose up --build -d

docker-down: ## stop docker compose services
	@docker-compose down

fmt: ## format go code
	@go fmt ./...

deps: ## download dependencies
	@go mod download
