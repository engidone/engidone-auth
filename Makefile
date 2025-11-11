# Engidone Auth Service - Makefile
# Microservicio de autenticación unificado con go-micro v5

# Variables
PROJECT_NAME := engidone-auth
VERSION := latest
SERVER_PORT := 8080
SERVICE_NAME := engidone-auth.service

# Paths
CMD_DIR := cmd
BIN_DIR := bin
SERVER_DIR := $(CMD_DIR)/server
CLIENT_DIR := $(CMD_DIR)/client
PROTO_DIR := proto
INTERNAL_DIR := internal

# Binaries
SERVER_BIN := $(BIN_DIR)/server
HEALTH_CLIENT_BIN := $(BIN_DIR)/health-client
SIGNIN_CLIENT_BIN := $(BIN_DIR)/signin-client
HELLO_CLIENT_BIN := $(BIN_DIR)/hello-client

# Go settings
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# Protocol Buffers settings
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTO_FILES := $(shell find $(INTERNAL_DIR) -name "*.proto")

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.ServiceName=$(SERVICE_NAME)"

# Help target
.PHONY: help
help: ## Show this help message
	@echo "=== $(PROJECT_NAME) Makefile ==="
	@echo ""
	@echo "Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make install-tools  # Install development tools"
	@echo "  make check          # Format and lint code"
	@echo "  make proto          # Generate protocol buffers"
	@echo "  make build          # Build all binaries"
	@echo "  make run-server     # Run the unified server"
	@echo "  make test-health    # Test health service"
	@echo "  make test-signin    # Test signin service"
	@echo "  make clean          # Clean build artifacts"

# Dependencies
.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install development tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "🛠️  Installing development tools..."
	@echo "Installing air for hot reloading..."
	@if ! command -v air >/dev/null 2>&1; then \
		$(GOGET) github.com/cosmtrek/air@latest; \
	else \
		echo "✅ air already installed"; \
	fi
	@echo "Installing goimports for import management..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		$(GOGET) golang.org/x/tools/cmd/goimports@latest; \
	else \
		echo "✅ goimports already installed"; \
	fi
	@echo "Installing golangci-lint for comprehensive linting..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$($(GOCMD) env GOPATH)/bin v1.54.2; \
	else \
		echo "✅ golangci-lint already installed"; \
	fi
	@echo "Installing additional linters..."
	@echo "  - golang.org/x/lint/golint (Go vet style linter)..."
	@if ! command -v golint >/dev/null 2>&1; then \
		$(GOGET) golang.org/x/lint/golint@latest; \
	else \
		echo "  ✅ golint already installed"; \
	fi
	@echo "  - honnef.co/go/tools/cmd/staticcheck (static analysis)..."
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		$(GOGET) honnef.co/go/tools/cmd/staticcheck@latest; \
	else \
		echo "  ✅ staticcheck already installed"; \
	fi
	@echo "  - github.com/golangci/go-misc/deadcode (dead code detection)..."
	@if ! command -v deadcode >/dev/null 2>&1; then \
		$(GOGET) github.com/golangci/go-misc/deadcode@latest; \
	else \
		echo "  ✅ deadcode already installed"; \
	fi
	@echo "Installing protocol buffer tools..."
	@if ! command -v $(PROTOC_GEN_GO) >/dev/null 2>&1; then \
		$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go; \
	else \
		echo "✅ protoc-gen-go already installed"; \
	fi
	@if ! command -v $(PROTOC_GEN_GO_GRPC) >/dev/null 2>&1; then \
		$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc; \
	else \
		echo "✅ protoc-gen-go-grpc already installed"; \
	fi
	@echo "✅ Development tools installation completed"

# Protocol Buffers generation
.PHONY: proto
proto: ## Generate protocol buffer files in ./proto
	@echo "🔧 Generating protocol buffers..."
	@mkdir -p $(PROTO_DIR)
	@if ! command -v $(PROTOC) >/dev/null 2>&1; then \
		echo "❌ protoc not found. Please install Protocol Buffers compiler."; \
		echo "   macOS: brew install protobuf"; \
		echo "   Ubuntu: sudo apt-get install protobuf-compiler"; \
		exit 1; \
	fi
	@if ! command -v $(PROTOC_GEN_GO) >/dev/null 2>&1; then \
		echo "📦 Installing protoc-gen-go..."; \
		$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go; \
	fi
	@if ! command -v $(PROTOC_GEN_GO_GRPC) >/dev/null 2>&1; then \
		echo "📦 Installing protoc-gen-go-grpc..."; \
		$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc; \
	fi
	@echo "📝 Generating Go files from .proto files..."
	@for proto_file in $(PROTO_FILES); do \
		echo "Processing $$proto_file..."; \
		$(PROTOC) \
			--go_out=$(PROTO_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(PROTO_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$proto_file; \
	done
	@echo "✅ Protocol buffers generated in $(PROTO_DIR)/"

.PHONY: proto-clean
proto-clean: ## Clean generated protocol buffer files
	@echo "🧹 Cleaning generated protocol buffers..."
	@rm -rf $(PROTO_DIR)
	@echo "✅ Protocol buffers cleaned"

# Formatting
.PHONY: fmt
fmt: ## Format Go code with goimports
	@echo "🎨 Formatting code with goimports..."
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w -local $(shell go mod edit -json | grep -o '"Module": "[^"]*"' | cut -d'"' -f4); \
		echo "✅ Code formatted with goimports"; \
	else \
		echo "⚠️  goimports not installed. Install with: make install-tools"; \
		echo "🔄 Falling back to gofmt..."; \
		$(GOFMT) -s -w .; \
	fi

# Linting
.PHONY: lint
lint: ## Lint Go code with comprehensive checks
	@echo "🔍 Running comprehensive linting..."
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo ""
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Checking import formatting with goimports..."; \
		unformatted=$$(find . -name "*.go" -not -path "./vendor/*" | xargs goimports -l -local $(shell go mod edit -json | grep -o '"Module": "[^"]*"' | cut -d'"' -f4)); \
		if [ -n "$$unformatted" ]; then \
			echo "❌ The following files need import formatting:"; \
			echo "$$unformatted"; \
			echo "💡 Run 'make fmt' to fix"; \
			exit 1; \
		else \
			echo "✅ Imports are properly formatted"; \
		fi; \
	else \
		echo "⚠️  goimports not installed. Install with: make install-tools"; \
	fi
	@echo ""
	@if command -v golint >/dev/null 2>&1; then \
		echo "Running golint..."; \
		lint_output=$$(golint ./... 2>/dev/null); \
		if [ -n "$$lint_output" ]; then \
			echo "⚠️  golint found issues:"; \
			echo "$$lint_output"; \
		else \
			echo "✅ golint passed"; \
		fi; \
	else \
		echo "⚠️  golint not installed. Install with: make install-tools"; \
	fi
	@echo ""
	@if command -v staticcheck >/dev/null 2>&1; then \
		echo "Running staticcheck..."; \
		staticcheck ./...; \
	else \
		echo "⚠️  staticcheck not installed. Install with: make install-tools"; \
	fi
	@echo ""
	@if command -v deadcode >/dev/null 2>&1; then \
		echo "Running deadcode detection..."; \
		deadcode ./...; \
	else \
		echo "⚠️  deadcode not installed. Install with: make install-tools"; \
	fi
	@echo ""
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: make install-tools"; \
	fi

# Quick format and lint check
.PHONY: check
check: fmt ## Format code and run comprehensive linting
	@echo "🔧 Running code quality checks..."
	$(MAKE) lint

# Testing
.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Test coverage
.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📈 Coverage report generated: coverage.html"

# Build targets
.PHONY: build
build: build-server build-clients ## Build all binaries

.PHONY: build-server
build-server: ## Build the unified server
	@echo "🏗️  Building unified server..."
	@mkdir -p $(BIN_DIR)
	cd $(SERVER_DIR) && $(GOBUILD) $(LDFLAGS) -o ../../$(SERVER_BIN) .

.PHONY: build-clients
build-clients: build-health-client build-signin-client build-hello-client ## Build all clients

.PHONY: build-health-client
build-health-client: ## Build health client
	@echo "🏥 Building health client..."
	@mkdir -p $(BIN_DIR)
	cd $(CLIENT_DIR)/health && $(GOBUILD) $(LDFLAGS) -o ../../../$(HEALTH_CLIENT_BIN) .

.PHONY: build-signin-client
build-signin-client: ## Build signin client
	@echo "🔐 Building signin client..."
	@mkdir -p $(BIN_DIR)
	cd $(CLIENT_DIR)/signin && $(GOBUILD) $(LDFLAGS) -o ../../../$(SIGNIN_CLIENT_BIN) .

.PHONY: build-hello-client
build-hello-client: ## Build hello client
	@echo "👋 Building hello client..."
	@mkdir -p $(BIN_DIR)
	cd $(CLIENT_DIR)/hello && $(GOBUILD) $(LDFLAGS) -o ../../../$(HELLO_CLIENT_BIN) .

# Run targets
.PHONY: run-server
run-server: ## Run the unified server
	@echo "🚀 Starting unified server..."
	@if [ ! -f $(SERVER_BIN) ]; then $(MAKE) build-server; fi
	@echo "📡 Server: $(SERVICE_NAME)"
	@echo "🌐 Port: $(SERVER_PORT)"
	@echo "🔧 Health Check & Signin Services"
	@echo ""
	$(SERVER_BIN)

.PHONY: run-health-client
run-health-client: ## Run health client
	@echo "🏥 Running health client..."
	@if [ ! -f $(HEALTH_CLIENT_BIN) ]; then $(MAKE) build-health-client; fi
	$(HEALTH_CLIENT_BIN)

.PHONY: run-signin-client
run-signin-client: ## Run signin client
	@echo "🔐 Running signin client..."
	@if [ ! -f $(SIGNIN_CLIENT_BIN) ]; then $(MAKE) build-signin-client; fi
	$(SIGNIN_CLIENT_BIN)

.PHONY: run-hello-client
run-hello-client: ## Run hello client
	@echo "👋 Running hello client..."
	@if [ ! -f $(HELLO_CLIENT_BIN) ]; then $(MAKE) build-hello-client; fi
	$(HELLO_CLIENT_BIN) $(filter-out $@,$(MAKECMDGOALS))

# Test targets (using clients)
.PHONY: test-health
test-health: build-health-client ## Test health service
	@echo "🏥 Testing Health Service..."
	@echo "⏳ Waiting for server to be ready..."
	@sleep 2
	$(HEALTH_CLIENT_BIN)

.PHONY: test-signin
test-signin: build-signin-client ## Test signin service
	@echo "🔐 Testing Signin Service..."
	@echo "⏳ Waiting for server to be ready..."
	@sleep 2
	$(SIGNIN_CLIENT_BIN)

.PHONY: test-hello
test-hello: build-hello-client ## Test hello service
	@echo "👋 Testing Hello Service..."
	@echo "⏳ Waiting for server to be ready..."
	@sleep 2
	$(HELLO_CLIENT_BIN) "Juan"

.PHONY: test-all
test-all: build-clients ## Test all services
	@echo "🧪 Testing all services..."
	@echo "⏳ Waiting for server to be ready..."
	@sleep 2
	@echo "=== Testing Health Service ==="
	$(HEALTH_CLIENT_BIN) &
	@sleep 1
	@echo "=== Testing Signin Service ==="
	$(SIGNIN_CLIENT_BIN)

# Development targets
.PHONY: dev
dev: ## Start development environment (server + test clients)
	@echo "🛠️  Starting development environment..."
	@echo "Starting server in background..."
	@$(MAKE) run-server &
	@SERVER_PID=$$!; \
	echo "Server PID: $$SERVER_PID"; \
	sleep 3; \
	echo "Running health tests..."; \
	$(MAKE) test-health; \
	sleep 2; \
	echo "Running signin tests..."; \
	$(MAKE) test-signin; \
	echo "Stopping server..."; \
	kill $$SERVER_PID 2>/dev/null || true

.PHONY: watch
watch: ## Watch for changes and rebuild with air hot reload
	@echo "👀 Watching for changes with air..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
	fi

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(PROJECT_NAME):$(VERSION) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p $(SERVER_PORT):$(SERVER_PORT) --name $(PROJECT_NAME) $(PROJECT_NAME):$(VERSION)

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	@echo "🛑 Stopping Docker container..."
	docker stop $(PROJECT_NAME) || true
	docker rm $(PROJECT_NAME) || true

# Clean targets
.PHONY: clean
clean: proto-clean ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.log

.PHONY: clean-all
clean-all: clean ## Clean everything including dependencies
	@echo "🧹 Cleaning everything..."
	$(GOMOD) clean -modcache
	@rm -f go.sum

# Status targets
.PHONY: status
status: ## Show project status
	@echo "=== $(PROJECT_NAME) Status ==="
	@echo "Version: $(VERSION)"
	@echo "Service: $(SERVICE_NAME)"
	@echo "Port: $(SERVER_PORT)"
	@echo ""
	@echo "📁 Binaries:"
	@if [ -f $(SERVER_BIN) ]; then echo "  ✅ Server: $(SERVER_BIN)"; else echo "  ❌ Server: not built"; fi
	@if [ -f $(HEALTH_CLIENT_BIN) ]; then echo "  ✅ Health Client: $(HEALTH_CLIENT_BIN)"; else echo "  ❌ Health Client: not built"; fi
	@if [ -f $(SIGNIN_CLIENT_BIN) ]; then echo "  ✅ Signin Client: $(SIGNIN_CLIENT_BIN)"; else echo "  ❌ Signin Client: not built"; fi
	@echo ""
	@echo "📊 Go Module:"
	@$(GOCMD) list -m

# Quick commands
.PHONY: quick-server
quick-server: build-server run-server ## Build and run server (quick command)

.PHONY: quick-test
quick-test: build-clients test-all ## Build and test all services (quick command)

# Install targets
.PHONY: install
install: build ## Install binaries to GOPATH/bin
	@echo "📦 Installing binaries..."
	@mkdir -p $$($(GOCMD) env GOPATH)/bin
	@cp $(SERVER_BIN) $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-server
	@cp $(HEALTH_CLIENT_BIN) $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-health-client
	@cp $(SIGNIN_CLIENT_BIN) $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-signin-client
	@echo "✅ Binaries installed to $$($(GOCMD) env GOPATH)/bin"

# Uninstall targets
.PHONY: uninstall
uninstall: ## Remove installed binaries
	@echo "🗑️  Uninstalling binaries..."
	@rm -f $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-server
	@rm -f $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-health-client
	@rm -f $$($(GOCMD) env GOPATH)/bin/$(PROJECT_NAME)-signin-client
	@echo "✅ Binaries uninstalled"

# Default target
.DEFAULT_GOAL := help

# Ensure all PHONY targets are declared
.PHONY: deps install-tools proto proto-clean fmt lint check test test-coverage build build-server build-clients build-health-client build-signin-client run-server run-health-client run-signin-client test-health test-signin test-all dev watch docker-build docker-run docker-stop clean clean-all status quick-server quick-test install uninstall