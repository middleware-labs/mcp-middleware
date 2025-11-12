.PHONY: build run test clean install lint fmt help

BINARY_NAME=mcp-middleware
BUILD_DIR=.
GO=go

build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

test:
	@echo "Running tests..."
	@$(GO) test -v ./test/...

test-coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./test/...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	@$(GO) test -race -v ./test/...

test-config:
	@echo "Running config tests..."
	@$(GO) test -v ./test/config

test-middleware:
	@echo "Running middleware tests..."
	@$(GO) test -v ./test/middleware

test-server:
	@echo "Running server tests..."
	@$(GO) test -v ./test/server

test-integration:
	@echo "Running integration tests..."
	@$(GO) test -v ./test/integration

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

install:
	@echo "Installing dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Dependencies installed"

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Format complete"

init-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo ".env file created. Please edit it with your credentials."; \
	else \
		echo ".env file already exists"; \
	fi

check-env:
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found!"; \
		echo "📝 Run: make init-env"; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | xargs) && \
	if [ -z "$$MIDDLEWARE_API_KEY" ] || [ -z "$$MIDDLEWARE_BASE_URL" ]; then \
		echo "❌ Error: Required environment variables not set!"; \
		echo "   MIDDLEWARE_API_KEY: $${MIDDLEWARE_API_KEY:-NOT SET}"; \
		echo "   MIDDLEWARE_BASE_URL: $${MIDDLEWARE_BASE_URL:-NOT SET}"; \
		echo ""; \
		echo "📝 Please edit .env file with your credentials"; \
		exit 1; \
	fi
	@echo "✓ Environment configuration is valid"

build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Multi-platform build complete"

stop-inspect:
	@echo "🛑 Stopping MCP Inspector..."
	@-pkill -f "modelcontextprotocol/inspector" 2>/dev/null && echo "✓ Inspector stopped" || echo "✓ No inspector running"

inspect: build check-env
	@echo "🔍 Starting MCP Inspector..."
	@-pkill -f "modelcontextprotocol/inspector" 2>/dev/null || true
	@bash -c '\
		set -a; \
		export $$(grep -v "^#" .env | grep -v "^$$" | xargs); \
		set +a; \
		echo "  API Key: $${MIDDLEWARE_API_KEY:0:10}..."; \
		echo "  Base URL: $$MIDDLEWARE_BASE_URL"; \
		echo ""; \
		echo "🚀 Launching MCP Inspector..."; \
		echo "   Inspector will open in your browser"; \
		echo "   Transport: stdio (command-line)"; \
		echo "   Press Ctrl+C to stop"; \
		echo ""; \
		exec npx --yes @modelcontextprotocol/inspector \
			--transport stdio \
			-e MIDDLEWARE_API_KEY="$$MIDDLEWARE_API_KEY" \
			-e MIDDLEWARE_BASE_URL="$$MIDDLEWARE_BASE_URL" \
			-e APP_MODE="$${APP_MODE:-stdio}" \
			./$(BINARY_NAME) \
	'

inspect-env: inspect

help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  test-config   - Run config tests only"
	@echo "  test-middleware - Run middleware tests only"
	@echo "  test-server   - Run server tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install dependencies"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  fmt           - Format code"
	@echo "  init-env      - Create .env file from .env.example"
	@echo "  check-env     - Validate environment configuration"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  inspect       - Run MCP Inspector for testing (requires Node.js/npx)"
	@echo "  stop-inspect  - Stop running MCP Inspector"
	@echo "  inspect-env   - Run MCP Inspector with .env configuration"
	@echo "  help          - Show this help message"

