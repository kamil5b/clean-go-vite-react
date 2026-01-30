# Setup Guide

## Quick Start

### Prerequisites
- Go 1.25+
- Node.js 18+
- Redis (for async task processing)

### Development Setup

1. **Install dependencies:**
```bash
make install-deps
```

2. **Start Redis (in another terminal):**
```bash
docker-compose up redis
```

3. **Run development mode:**
```bash
make dev
```

This starts:
- Go server on `http://localhost:8080`
- Frontend dev server on `http://localhost:5173`
- API accessible at `http://localhost:8080/api`

### Running Components Separately

**Server only:**
```bash
make server
```

**Worker only:**
```bash
make worker
```

Both require Redis to be running.

## Architecture Overview

The codebase follows a clean layered architecture:

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access)
                      ↓
                  Task (Async)
                      ↓
                  Worker (Processing)
```

### Key Directories

- `cmd/server/` - HTTP server entry point
- `cmd/worker/` - Background worker entry point
- `internal/api/` - HTTP handlers and routing
- `internal/service/` - Business logic (transport-agnostic)
- `internal/repository/` - Data access interfaces
- `internal/task/` - Async task definitions
- `internal/worker/` - Task processors
- `internal/model/` - Domain models
- `internal/platform/` - Infrastructure setup
- `internal/web/` - Frontend hosting (optional)

## Configuration

Environment variables (see `env.example`):

- `SERVER_PORT` - Server port (default: 8080)
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `ASYNQ_ENABLED` - Enable async processing (default: true)
- `DEV_MODE` - Enable development mode (default: false)

Create a `.env` file in the root directory to override defaults.

## Production Build

```bash
make build
```

Produces:
- `bin/server` - HTTP server binary
- `bin/worker` - Worker binary
- `frontend/dist/` - Built frontend assets

## Testing

```bash
make test
```

Runs all unit tests across the project.

## Adding New Features

### Adding a New API Endpoint

1. Create service method in `internal/service/`
2. Create handler in `internal/api/handler/`
3. Register route in `internal/api/router.go`

### Adding Async Tasks

1. Define task in `internal/task/tasks.go`
2. Create processor in `internal/worker/`
3. Register processor in `cmd/worker/main.go`
4. Enqueue from handler using Asynq client

## Deployment

Both server and worker can be deployed separately:

```bash
# Deploy server
./bin/server

# Deploy worker (requires Redis)
./bin/worker
```

Configure via environment variables.

## Troubleshooting

**Redis connection failed:**
- Ensure Redis is running: `docker-compose up redis`
- Check `REDIS_HOST` and `REDIS_PORT` settings

**Frontend not loading:**
- Ensure `DEV_MODE=true` is set for dev mode
- In production, run `make build` first

**Hot reload not working:**
- Frontend dev server should be running on port 5173
- Check browser console for errors