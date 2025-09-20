.PHONY: help dev dev-stop dev-reset run stop clean

# Default target
help:
	@echo "Available commands:"
	@echo "  dev       - Start development environment (DB + server with seeding)"
	@echo "  dev-stop  - Stop development environment"
	@echo "  dev-reset - Reset database and restart dev environment"
	@echo "  run       - Run server only (assumes DB is running)"
	@echo "  stop      - Stop server only"
	@echo "  clean     - Clean up all containers and volumes"

# Development environment with database
dev: dev-stop
	@echo "Starting development environment..."
	@cd .. && docker-compose up -d db
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Starting server with seeding..."
	@DEV_MODE=true DATABASE_URL="postgres://postgres:postgres@localhost:5432/pin?sslmode=disable" go run cmd/server/main.go

# Stop development environment
dev-stop:
	@echo "Stopping development environment..."
	@cd .. && docker-compose down
	@pkill -f "go run cmd/server/main.go" 2>/dev/null || true

# Reset database and restart
dev-reset: dev-stop
	@echo "Resetting database..."
	@cd .. && docker volume rm pin_pg_data 2>/dev/null || true
	@cd .. && docker-compose up -d db
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Starting server with fresh seeding..."
	@DEV_MODE=true DATABASE_URL="postgres://postgres:postgres@localhost:5432/pin?sslmode=disable" go run cmd/server/main.go

# Run server only (assumes DB is running)
run:
	@echo "Starting server..."
	@DEV_MODE=true DATABASE_URL="postgres://postgres:postgres@localhost:5432/pin?sslmode=disable" go run cmd/server/main.go

# Stop server only
stop:
	@pkill -f "go run cmd/server/main.go" 2>/dev/null || true

# Clean up everything
clean: dev-stop
	@echo "Cleaning up containers and volumes..."
	@cd .. && docker-compose down --volumes --remove-orphans
	@docker volume rm pin_pg_data 2>/dev/null || true