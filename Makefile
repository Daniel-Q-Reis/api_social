# Makefile for Social Media API

# Variables
APP_NAME=social-api
GO_BUILD=go build
GO_RUN=go run
GO_TEST=go test
GO_MOD=go mod
GO_CLEAN=go clean

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	$(GO_BUILD) -o $(APP_NAME) cmd/app/main.go

# Run the application
.PHONY: run
run:
	$(GO_RUN) cmd/app/main.go

# Install dependencies
.PHONY: deps
deps:
	$(GO_MOD) tidy

# Clean build files
.PHONY: clean
clean:
	$(GO_CLEAN)
	rm -f $(APP_NAME)

# Run tests
.PHONY: test
test:
	$(GO_TEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO_TEST) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Install linter
.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all checks
.PHONY: check
check: fmt vet lint test

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all             - Build the application (default)"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  deps            - Install dependencies"
	@echo "  clean           - Clean build files"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  fmt             - Format code"
	@echo "  vet             - Vet code"
	@echo "  lint            - Run linter"
	@echo "  lint-install    - Install linter"
	@echo "  check           - Run all checks"
	@echo "  help            - Show this help"

# Integration test target
.PHONY: compose-up-integration-test
compose-up-integration-test:
	docker-compose -f docker-compose-integration-test.yml up -d --build
	go test -v ./integration-test/...
	docker-compose -f docker-compose-integration-test.yml down