MOCKGEN := mockgen
REPO_INTERFACES_DIR := backend/repository/interfaces
REPO_MOCK_DIR := backend/repository/mock

# Detect OS and Architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
BINARY_NAME := server
BINARY_PATH := ./bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH)
ifeq ($(GOOS),windows)
	BINARY_PATH := ./bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe
endif

.PHONY: help dev server build test install-deps clean repository-mocks

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install dependencies"
	@echo "  make dev           - Start development (frontend + server)"
	@echo "  make server        - Start HTTP server only"
	@echo "  make build         - Build production binary for current OS/Arch"
	@echo "  make build-all     - Build production binaries for all platforms"
	@echo "  make build-linux   - Build for Linux (amd64)"
	@echo "  make build-windows - Build for Windows (amd64)"
	@echo "  make build-darwin  - Build for macOS (amd64)"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Clean build artifacts"

install-deps:
	go mod tidy
	go mod download
	cd frontend && yarn install

dev:
	@echo "Starting development environment..."
	cd frontend && yarn dev & sleep 1 && DEV_MODE=true air

server:
	DEV_MODE=true air

build:
	@echo "Building frontend..."
	cd frontend && yarn build
	@echo "Building server binary for $(GOOS)/$(GOARCH)..."
	@mkdir -p ./bin
	ENV=prod GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildvcs=false -o $(BINARY_PATH) ./cmd/server/main.go
	@echo "Binary created at: $(BINARY_PATH)"

build-all: build-linux build-windows build-darwin build-linux-arm64 build-darwin-arm64
	@echo "All binaries built successfully"

build-linux:
	@echo "Building for Linux (amd64)..."
	cd frontend && yarn build
	@mkdir -p ./bin
	ENV=prod GOOS=linux GOARCH=amd64 go build -buildvcs=false -o ./bin/$(BINARY_NAME)-linux-amd64 ./cmd/server/main.go

build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	cd frontend && yarn build
	@mkdir -p ./bin
	ENV=prod GOOS=linux GOARCH=arm64 go build -buildvcs=false -o ./bin/$(BINARY_NAME)-linux-arm64 ./cmd/server/main.go

build-windows:
	@echo "Building for Windows (amd64)..."
	cd frontend && yarn build
	@mkdir -p ./bin
	ENV=prod GOOS=windows GOARCH=amd64 go build -buildvcs=false -o ./bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/server/main.go

build-darwin:
	@echo "Building for macOS (amd64)..."
	cd frontend && yarn build
	@mkdir -p ./bin
	ENV=prod GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o ./bin/$(BINARY_NAME)-darwin-amd64 ./cmd/server/main.go

build-darwin-arm64:
	@echo "Building for macOS (arm64/M1)..."
	cd frontend && yarn build
	@mkdir -p ./bin
	ENV=prod GOOS=darwin GOARCH=arm64 go build -buildvcs=false -o ./bin/$(BINARY_NAME)-darwin-arm64 ./cmd/server/main.go

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

repository-mocks:
	@rm -rf $(REPO_MOCK_DIR)
	@mkdir -p $(REPO_MOCK_DIR)
	@for f in $(REPO_INTERFACES_DIR)/*.go; do \
		base=$$(basename $$f .go); \
		name=$${base%_interface}; \
		echo "Generating mock for interface: $$name"; \
		$(MOCKGEN) \
			-source=$$f \
			-destination=$(REPO_MOCK_DIR)/$${name}_mock.go \
			-package=mock; \
	done
