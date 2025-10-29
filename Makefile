.PHONY: help build test test-verbose test-coverage lint fmt clean run benchmark install-tools

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=customer-importer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	@echo "Binary built: $(BINARY_NAME)"

run: ## Run the application with default settings
	$(GOCMD) run main.go

run-verbose: ## Run the application with verbose logging
	$(GOCMD) run main.go -verbose

test: ## Run all tests
	$(GOTEST) -v ./...

test-verbose: ## Run all tests with verbose output
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage report
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@$(GOCMD) tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run 'make install-tools' first." && exit 1)
	golangci-lint run ./...

fmt: ## Format all Go files
	$(GOFMT) -s -w .
	@which goimports > /dev/null && goimports -w . || echo "goimports not installed, skipping"

vet: ## Run go vet
	$(GOCMD) vet ./...

tidy: ## Tidy and verify go.mod
	$(GOMOD) tidy
	$(GOMOD) verify

clean: ## Remove build artifacts and coverage files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -f customerimporter/test_output.csv
	rm -f exporter/test_output.csv

install-tools: ## Install development tools (golangci-lint, goimports)
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goimports..."
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed successfully"

ci: fmt lint vet test-coverage ## Run all CI checks (format, lint, vet, test with coverage)
	@echo "All CI checks passed!"

all: clean fmt lint vet test build ## Run clean, format, lint, vet, test, and build
	@echo "Build complete!"
