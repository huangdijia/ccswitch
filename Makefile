.PHONY: build install clean test help

# Binary name
BINARY_NAME=ccswitch

# Installation paths
PREFIX?=/usr/local
BINDIR=$(PREFIX)/bin
SHAREDIR=$(PREFIX)/share/$(BINARY_NAME)

# Build flags
GO=go

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u +'%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME) with version $(VERSION) (commit: $(COMMIT), date: $(DATE))"
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

install: build ## Install the binary and config files
	@echo "Installing $(BINARY_NAME) to $(BINDIR)..."
	@mkdir -p $(BINDIR)
	@mkdir -p $(SHAREDIR)
	@cp $(BINARY_NAME) $(BINDIR)/
	@chmod +x $(BINDIR)/$(BINARY_NAME)
	@cp -r config $(SHAREDIR)/
	@echo "Installation complete!"
	@echo "You can now run: $(BINARY_NAME) init"

uninstall: ## Uninstall the binary and config files
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(BINDIR)/$(BINARY_NAME)
	@rm -rf $(SHAREDIR)
	@echo "Uninstallation complete!"

clean: ## Remove build artifacts
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

test: ## Run tests
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt: ## Format the code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

mod: ## Download and tidy dependencies
	$(GO) mod download
	$(GO) mod tidy

all: clean build ## Clean and build

.DEFAULT_GOAL := help
