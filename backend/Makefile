.PHONY: run build test clean deps health-check

all: deps build

deps:
	go mod tidy
	go mod download

build:
	go build -o bin/server cmd/server/main.go

PORT ?= 8080
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= pin
DB_HOST ?= localhost
DB_PORT ?= 5432
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

run:
	PORT=$(PORT) DATABASE_URL=$(DATABASE_URL) go run cmd/server/main.go

db-up:
	docker compose up -d db

db-down:
	docker compose down

migrate-up:
	DATABASE_URL=$(DATABASE_URL) docker compose run --rm migrate

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=add_users_table" && exit 1)
	@mkdir -p migrations
	docker run --rm -v $(PWD)/migrations:/migrations -w /migrations migrate/migrate:v4.17.1 create -ext sql -seq $(name)

test:
	go test ./...

clean:
	rm -rf bin/

install:
	go mod tidy
