# Todo API Makefile

.PHONY: help run build test migrate clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make run       - Run the API server"
	@echo "  make build     - Build the binary"
	@echo "  make migrate   - Run all database migrations"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deps      - Download dependencies"

# Run the API server
run:
	go run cmd/api/main.go

# Build the binary
build:
	go build -o bin/todo-api cmd/api/main.go

# Run all database migrations
migrate:
	@echo "Running database migrations..."
	@./scripts/migrate.sh

# Download dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean
