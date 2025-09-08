.PHONY: run build test clean deps health-check

all: deps build

deps:
	go mod tidy
	go mod download

build:
	go build -o bin/server cmd/server/main.go

PORT ?= 8080

run:
	PORT=$(PORT) go run cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/

install:
	go mod tidy
