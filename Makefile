# Variables
PROJECT_NAME := wtf
MAIN_PATH := main.go
OUTPUT_DIR := bin
OUTPUT_BIN := $(OUTPUT_DIR)/$(PROJECT_NAME)
GO := go
VERSION ?= dev

# Default target
help:
	@echo "$(PROJECT_NAME) - Infrastructure as Code with Terraform JSON Generation"
	@echo ""
	@echo "Available targets:"
	@echo "  build          Build the binary"
	@echo "  test           Run all tests"
	@echo "  test-verbose   Run tests with verbose output"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  install        Build and install the binary to \$$GOPATH/bin"
	@echo "  run            Run the binary directly"
	@echo "  clean          Remove build artifacts"
	@echo "  lint           Run linter (requires golangci-lint)"
	@echo "  fmt            Format code with gofmt"
	@echo "  vet            Run go vet"
	@echo "  mod-tidy       Tidy go.mod dependencies"
	@echo "  coverage       Generate coverage report (HTML)"
	@echo "  release        Build release binaries for multiple platforms"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make install"
	@echo "  make run -- --help"

# ==============================================================================
# BUILD TARGETS
# ==============================================================================

## Build the binary for current platform
build: fmt vet
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -o $(OUTPUT_BIN) $(MAIN_PATH)
	@echo "✓ Binary built: $(OUTPUT_BIN)"

## Build for Linux
build-linux:
	@echo "Building $(PROJECT_NAME) for Linux..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -o $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "✓ Built: $(OUTPUT_DIR)/$(PROJECT_NAME)-linux-amd64"

## Build for macOS
build-macos:
	@echo "Building $(PROJECT_NAME) for macOS..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -o $(OUTPUT_DIR)/$(PROJECT_NAME)-macos-amd64 $(MAIN_PATH)
	@echo "✓ Built: $(OUTPUT_DIR)/$(PROJECT_NAME)-macos-amd64"

## Build for Windows
build-windows:
	@echo "Building $(PROJECT_NAME) for Windows..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -o $(OUTPUT_DIR)/$(PROJECT_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "✓ Built: $(OUTPUT_DIR)/$(PROJECT_NAME)-windows-amd64.exe"

## Build release binaries for all platforms
release: clean build-linux build-macos build-windows
	@echo "✓ All release binaries built"
	@ls -lh $(OUTPUT_DIR)/$(PROJECT_NAME)-*

# ==============================================================================
# TEST TARGETS
# ==============================================================================

## Run all tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v -race

## Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GO) test ./... -v -race -count=1

## Run tests with coverage
test-coverage: coverage
	@echo "✓ Coverage report generated: coverage.html"

## Generate coverage report (HTML)
coverage:
	@echo "Generating coverage report..."
	$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

## Run E2E tests only
test-e2e:
	@echo "Running E2E tests..."
	$(GO) test ./tests/e2e -v

## Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GO) test ./internal/... -v

# ==============================================================================
# INSTALL TARGETS
# ==============================================================================

## Build and install binary to $$GOPATH/bin
install: build
	@echo "Installing $(PROJECT_NAME)..."
	$(GO) install $(MAIN_PATH)
	@echo "✓ Installed to $$($(GO) env GOPATH)/bin/$(PROJECT_NAME)"

## Install dependencies
install-deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	@echo "✓ Dependencies installed"

## Install development tools
install-dev-tools:
	@echo "Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Development tools installed"

# ==============================================================================
# CODE QUALITY TARGETS
# ==============================================================================

## Format code with gofmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Code formatted"

## Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "✓ No issues found"

## Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: make install-dev-tools" && exit 1)
	golangci-lint run ./...
	@echo "✓ Linting complete"

## Run all checks (fmt, vet, lint)
check: fmt vet lint
	@echo "✓ All checks passed"

# ==============================================================================
# UTILITY TARGETS
# ==============================================================================

## Run the binary
run: build
	@$(OUTPUT_BIN)

## Run with arguments (e.g., make run -- --help)
run-args: build
	@$(OUTPUT_BIN) $(filter-out $@,$(MAKECMDGOALS))

## Tidy go.mod
mod-tidy:
	@echo "Tidying go.mod..."
	$(GO) mod tidy
	@echo "✓ go.mod tidied"

## Verify go.mod
mod-verify:
	@echo "Verifying go.mod..."
	$(GO) mod verify
	@echo "✓ go.mod verified"

## Download all dependencies
mod-download:
	@echo "Downloading dependencies..."
	$(GO) mod download -x
	@echo "✓ Dependencies downloaded"

## Update all dependencies
mod-update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✓ Dependencies updated"

## Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out coverage.html
	@$(GO) clean
	@echo "✓ Cleaned"

## Show version info
version:
	@echo "$(PROJECT_NAME) Build Info:"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Go Version: $$($(GO) version)"

# ==============================================================================
# CI/CD TARGETS
# ==============================================================================

## Run full CI pipeline (check, test, coverage, build)
ci: check test coverage build
	@echo "✓ CI pipeline complete"

## Run tests and build (lightweight CI)
ci-lite: test build
	@echo "✓ CI-lite pipeline complete"

## Pre-commit checks
pre-commit: fmt vet test-unit
	@echo "✓ Pre-commit checks passed"

# ==============================================================================
# DOCUMENTATION TARGETS
# ==============================================================================

## Show project structure
structure:
	@echo "Project structure:"
	@tree -L 2 -I 'bin|*.o|*.a' 2>/dev/null || find . -maxdepth 2 -type f -name "*.go" | head -20

## Show available Go commands
go-help:
	@$(GO) help

## Show go test help
test-help:
	@$(GO) test -help

.PHONY: help build build-linux build-macos build-windows release test test-verbose test-coverage test-e2e test-unit coverage install install-deps install-dev-tools fmt vet lint check run run-args mod-tidy mod-verify mod-download mod-update clean version ci ci-lite pre-commit structure go-help test-help
