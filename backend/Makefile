.PHONY: build run dev test lint tidy clean swagger deps ent/gen ent/new migrate/diff migrate/apply

# Build variables
BINARY_NAME=server
BUILD_DIR=bin

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
AIR_BIN=$(shell command -v air 2>/dev/null || echo "$(HOME)/go/bin/air")

# Load environment variables from .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Build the server
build:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "✓ Server built: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the server
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@if [ -x "$(AIR_BIN)" ]; then \
		"$(AIR_BIN)" -c .air.toml; \
	else \
		echo "air not found; running without hot reload"; \
		echo "Install with: go install github.com/air-verse/air@v1.61.7"; \
		$(GORUN) ./cmd/server; \
	fi

# Run tests
test:
	@$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Run linter (requires golangci-lint)
lint:
	@golangci-lint run ./...

# Tidy dependencies
tidy:
	@$(GOMOD) tidy

# Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

# Generate swagger docs (requires swag)
swagger:
	@echo "Generating swagger docs..."
	@swag init -g cmd/server/main.go --parseDependency --parseInternal
	@echo "✓ Swagger docs generated in ./docs"

# Install dependencies (tools)
deps:
	@echo "Installing dependencies..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Dependencies installed"

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the server binary"
	@echo "  run            - Build and run the server"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  tidy           - Tidy dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  swagger        - Generate swagger docs"
	@echo "  ent/gen        - Generate Ent code"
	@echo "  ent/new        - Create a new Ent schema (usage: make ent/new NAME=User)"
	@echo "  migrate/diff   - Create a migration file based on schema changes (usage: make migrate/diff NAME=add_users)"
	@echo "  migrate/apply  - Apply pending migrations to the database (usage: make migrate/apply DB_URL=...)"

# Generate Ent code
ent/gen:
	@go generate ./internal/ent

# Create a new Ent schema (usage: make ent/new NAME=User)
ent/new:
	@go run -mod=mod entgo.io/ent/cmd/ent new --target internal/ent/schema $(NAME)

# Create a new migration file based on Ent schema changes
migrate/diff:
	@atlas migrate diff $(NAME) \
	  --dir "file://internal/ent/migrate/migrations" \
	  --to "ent://internal/ent/schema" \
	  --dev-url "docker://postgres/15/test"

# Database migration variables
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= go_backend_template
DB_SSLMODE ?= disable
DB_URL ?= "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

migrate/apply:
	@atlas migrate apply \
	  --dir "file://internal/ent/migrate/migrations" \
	  --url $(DB_URL)
