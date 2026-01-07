# ==========================================
# pendek-in - URL Shortener Service
# ==========================================

# Binary name
BINARY_NAME=pendek-in
MAIN_PATH=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

# Database migration paths
MIGRATIONS_DIR=internal/database/schema

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Build DATABASE_URL from .env values
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= link-short
DB_SSLMODE ?= disable
DATABASE_URL = host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)

.PHONY: all build run clean test fmt vet lint deps tidy sqlc swagger migrate-up migrate-down migrate-status migrate-create help

# Default target
all: build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

# ==========================================
# Build & Run
# ==========================================

## build: Build the application binary
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)

## run: Run the application
run:
	$(GORUN) $(MAIN_PATH)

## run-watch: Run with hot reload (requires air)
run-watch:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

## clean: Remove build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/

# ==========================================
# Dependencies
# ==========================================

## deps: Download dependencies
deps:
	$(GOMOD) download

## tidy: Tidy and verify dependencies
tidy:
	$(GOMOD) tidy
	$(GOMOD) verify

## vendor: Vendor dependencies
vendor:
	$(GOMOD) vendor

# ==========================================
# Code Quality
# ==========================================

## fmt: Format code
fmt:
	$(GOFMT) -s -w .

## vet: Run go vet
vet:
	$(GOVET) ./...

## lint: Run golangci-lint (install if not present)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## test: Run tests
test:
	$(GOTEST) -v -race -cover ./...

## test-short: Run short tests
test-short:
	$(GOTEST) -v -short ./...

## coverage: Run tests with coverage report
coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# ==========================================
# Database & SQL Generation
# ==========================================

## sqlc: Generate Go code from SQL (requires sqlc)
sqlc:
	@which sqlc > /dev/null || (echo "Installing sqlc..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
	sqlc generate

## sqlc-verify: Verify sqlc configuration
sqlc-verify:
	sqlc compile

## swagger: Generate Swagger documentation (requires swag)
swagger:
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init

## swagger-fmt: Format Swagger comments
swagger-fmt:
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag fmt

## migrate-up: Run all database migrations (requires goose)
migrate-up:
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

## migrate-down: Rollback the last migration
migrate-down:
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

## migrate-status: Show migration status
migrate-status:
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

## migrate-reset: Rollback all migrations
migrate-reset:
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset

## migrate-create: Create a new migration file (usage: make migrate-create name=migration_name)
migrate-create:
	@which goose > /dev/null || (echo "Installing goose..." && go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir $(MIGRATIONS_DIR) create $(name) sql
	
# ==========================================
# Development Setup
# ==========================================

## setup: Install development tools
setup:
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "All tools installed!"

## env: Copy .env.example to .env (if .env.example exists)
env:
	@if [ -f .env.example ] && [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	elif [ -f .env ]; then \
		echo ".env file already exists"; \
	else \
		echo "No .env.example file found"; \
	fi

# ==========================================
# CI/CD Helpers
# ==========================================

## ci: Run all CI checks
ci: fmt vet lint test

## build-all: Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
