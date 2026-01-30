.PHONY: help dev server worker build test install-deps clean

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install dependencies"
	@echo "  make dev           - Start development (frontend + server)"
	@echo "  make server        - Start HTTP server only"
	@echo "  make worker        - Start async worker"
	@echo "  make build         - Build production binaries"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Clean build artifacts"

install-deps:
	go mod tidy
	go mod download
	cd frontend && yarn install

dev:
	@echo "Starting development environment..."
	cd frontend && yarn dev & sleep 3 && DEV_MODE=true go run ./cmd/server

server:
	DEV_MODE=true go run ./cmd/server

worker:
	go run ./cmd/worker

build:
	@echo "Building frontend..."
	cd frontend && yarn build
	@echo "Building server binary..."
	ENV=prod go build -buildvcs=false -o ./bin/server ./cmd/server/main.go
	@echo "Building worker binary..."
	ENV=prod go build -buildvcs=false -o ./bin/worker ./cmd/worker/main.go

test:
	go test -v -cover -race ./...

test-verbose:
	go test -v -cover -race -failfast ./...

test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf ./bin
	rm -f coverage.out coverage.html
	go clean -testcache
