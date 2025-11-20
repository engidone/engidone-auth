.PHONY: build run dev-setup help clean docker-build docker-up docker-down docker-reset sqlc migrate migrate-up migrate-create migrate-rollback migrate-to migrate-rollback-to test lint format check-deps

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building application..."
	CGO_ENABLED=1 go build -o bin/server ./cmd/server

# Run the application
run: build
	@echo "🚀 Running application..."
	./bin/server

# Run with hot reload (requires air)
dev:
	@command -v air >/dev/null 2>&1 || { echo "❌ 'air' not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1; }
	air

# Development setup (install deps, generate code, build)
dev-setup:
	@echo "⚙️ Setting up development environment..."
	go mod tidy
	go mod download
	$(MAKE) sqlc
	$(MAKE) build

# Run database migrations (alias for migrate-up)
migrate: migrate-up

# Apply all pending migrations
migrate-up:
	@echo "🗄️ Running all pending migrations..."
	@if [ -d "./migrations" ]; then \
		go run cmd/db/migrate.go; \
	else \
		echo "❌ No migrations directory found at ./migrations"; \
	fi

# Create new migration
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "📄 Creating new migration: $(name)..."
	@if [ -d "./migrations" ]; then \
		go run github.com/pressly/goose/v3/cmd/goose -s -dir ./migrations create $(name) sql; \
	else \
		echo "❌ No migrations directory found at ./migrations"; \
	fi

# Rollback migrations (default 1 step)
migrate-rollback:
	@if [ -z "$(steps)" ]; then \
		echo "🗄️ Rolling back 1 migration..."; \
		go run cmd/db/migrate.go -r -k 1; \
	else \
		echo "🗄️ Rolling back $(steps) migration(s)..."; \
		go run cmd/db/migrate.go -r -k $(steps); \
	fi

# Migrate to specific sequential
migrate-to:
	@if [ -z "$(seq)" ]; then \
		echo "❌ Usage: make migrate-to seq=5"; \
		exit 1; \
	fi
	@echo "🗄️ Migrating to $(seq)..."
	@if [ -d "./migrations" ]; then \
		go run cmd/db/migrate.go -s $(seq); \
	else \
		echo "❌ No migrations directory found at ./migrations"; \
	fi

# Rollback to specific sequential
migrate-rollback-to:
	@if [ -z "$(seq)" ]; then \
		echo "❌ Usage: make migrate-rollback-to seq=3"; \
		exit 1; \
	fi
	@echo "🗄️ Rolling back to $(seq)..."
	@if [ -d "./migrations" ]; then \
		go run cmd/db/migrate.go -r -s $(seq); \
	else \
		echo "❌ No migrations directory found at ./migrations"; \
	fi

# Generate SQLC code
sqlc:
	@echo "🔧 Generating SQLC code..."
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint code
lint:
	@echo "🔍 Linting code..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "❌ 'golangci-lint' not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run

# Format code
format:
	@echo "📝 Formatting code..."
	go fmt ./...
	goimports -w .

# Check for outdated dependencies
check-deps:
	@echo "📦 Checking for outdated dependencies..."
	go list -u -m -json all | go-mod-outdated -update -direct

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/ coverage.out coverage.html
	go clean

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t engidoneauth .

docker-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker Compose services..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-reset:
	@echo "🐳 Resetting Docker environment..."
	docker-compose down -v
	docker system prune -f
	$(MAKE) docker-up

# Help
help:
	@echo "🚀 EngiDone Auth - Available targets:"
	@echo ""
	@echo "📦 Build & Run:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  dev-setup     - Initial development setup"
	@echo ""
	@echo "🗄️ Database:"
	@echo "  migrate             - Run database migrations (alias for migrate-up)"
	@echo "  migrate-up          - Apply all pending migrations"
	@echo "  migrate-create      - Create new migration (usage: make migrate-create name=migration_name)"
	@echo "  migrate-rollback    - Rollback migrations (usage: make migrate-rollback steps=3)"
	@echo "  migrate-to          - Migrate to specific sequential (usage: make migrate-to seq=5)"
	@echo "  migrate-rollback-to - Rollback to specific sequential (usage: make migrate-rollback-to seq=3)"
	@echo "  sqlc                - Generate SQLC code"
	@echo ""
	@echo "🧪 Quality:"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  format        - Format code (requires goimports)"
	@echo "  check-deps    - Check for outdated dependencies"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop Docker Compose services"
	@echo "  docker-reset  - Reset Docker environment (removes volumes)"
	@echo "  docker-logs   - Show Docker logs"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  clean         - Clean build artifacts and temporary files"
	@echo "  help          - Show this help message"

# Project structure info
info:
	@echo "=== Package-Driven Design Structure ==="
	@echo "internal/"
	@echo "├── auth/           # Authentication domain"
	@echo "├── greet/          # Greeting domain"
	@echo "├── storage/        # Database layer"
	@echo "├── server/         # gRPC server"
	@echo "└── di/             # Dependency injection"
	@echo ""
	@echo "📖 For detailed structure see: internal/README.md"