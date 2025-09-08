# Pin Backend Makefile

.PHONY: run build test clean deps health-check

# Default target
all: deps build

# Download dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Health check (requires server to be running)
health-check:
	@echo "Checking server health..."
	@curl -s http://localhost:8080/health | jq . || curl -s http://localhost:8080/health

# Install dependencies
install:
	go mod tidy
