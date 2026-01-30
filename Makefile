.PHONY: dev build server worker clean help

help:
	@echo "Available commands:"
	@echo "  make dev           - Run development mode (frontend + server)"
	@echo "  make server        - Run server only"
	@echo "  make worker        - Run worker only"
	@echo "  make build         - Build production binary"
	@echo "  make clean         - Clean build artifacts"

dev:
	@echo "Starting development environment..."
	@DEV_MODE=true go run ./cmd/server &
	@cd frontend && yarn dev

server:
	@echo "Starting server..."
	@go run ./cmd/server

worker:
	@echo "Starting worker..."
	@go run ./cmd/worker

build:
	@echo "Building frontend..."
	@cd frontend && yarn build
	@echo "Building server binary..."
	@go build -buildvcs=false -o ./bin/server ./cmd/server
	@echo "Building worker binary..."
	@go build -buildvcs=false -o ./bin/worker ./cmd/worker
	@echo "Build complete!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@rm -rf ./frontend/dist
	@echo "Clean complete!"

install-deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@go mod tidy
	@echo "Installing frontend dependencies..."
	@cd frontend && yarn install
	@echo "All dependencies installed!"

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...
