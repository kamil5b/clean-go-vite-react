MOCKGEN := mockgen
REPO_INTERFACES_DIR := backend/repository/interfaces
REPO_MOCK_DIR := backend/repository/mock

.PHONY: help dev server build test install-deps clean repository-mocks

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install dependencies"
	@echo "  make dev           - Start development (frontend + server)"
	@echo "  make server        - Start HTTP server only"
	@echo "  make build         - Build production binaries"
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
	@echo "Building server binary..."
	ENV=prod go build -buildvcs=false -o ./bin/server ./cmd/server/main.go

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
```

Now let me update the docker-compose files and other files:
