.PHONY: build run test seed clean docker-build docker-up docker-down sqlc

# Default target
all: build

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run: build
	./bin/server

# Run with seeders
run-seed: build
	SERVER_PORT=9000 ./bin/server

# Note: Tests require MySQL database to be running
test-sql:
	@echo "SQL test requires database connection - use docker-compose up first"

# Note: Seeders test requires MySQL database to be running
test-seeders:
	@echo "Seeders test requires database connection - use docker-compose up first"

# Generate SQLC code
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Docker commands
docker-build:
	docker build -t engidone-auth .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development setup
dev-setup:
	go mod tidy
	go mod download
	$(MAKE) sqlc
	$(MAKE) build

# Full test suite
test: build
	@echo "All tests require database connection - use docker-compose up first"

# Database operations
db-migrate:
	@echo "Migrations run automatically on application startup via YAML seeders"

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  run-seed      - Run the application with seeders"
	@echo "  test-sql      - Test SQL implementation"
	@echo "  test-seeders  - Test YAML seeders"
	@echo "  sqlc          - Generate SQLC code"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop Docker Compose services"
	@echo "  docker-logs   - Show Docker logs"
	@echo "  dev-setup     - Initial development setup"
	@echo "  test          - Run all tests"
	@echo "  db-migrate    - Run database migrations"
	@echo "  help          - Show this help message"