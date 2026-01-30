# Architecture Overview

This document describes the complete architecture of the repository: runtime boundaries, dependency rules, backend structure (server + worker), asynchronous processing with Asynq, and frontend architecture. The primary goals are **clean separation**, **easy detachment**, and **long-term maintainability**.

---

## High-Level Goals

* Frontend and backend are loosely coupled
* Backend supports multiple runtimes (HTTP server, async worker)
* Business logic is shared and transport-agnostic
* Asynchronous work is handled outside the request lifecycle
* Future changes (detach frontend, swap queue, add gRPC) require minimal refactoring

---

## System Diagram

```
Frontend (Vite + React)
        │
        │ HTTP (/api)
        ▼
Server (Echo :8080)
        │
        │ Enqueue (Asynq client)
        ▼
Redis (Asynq)
        │
        ▼
Worker (Asynq server)
```

---

## Repository Layout

```
.
├── cmd/
│   ├── server/              # HTTP runtime (API + optional frontend hosting)
│   │   └── main.go
│   └── worker/              # Background worker runtime
│       └── main.go
│
├── internal/
│   ├── api/                 # HTTP transport layer
│   │   ├── handler/
│   │   └── router.go
│   │
│   ├── service/             # Business logic (shared)
│   │
│   ├── repository/          # Data access layer
│   │   ├── interfaces
│   │   └── implementations
│   │
│   ├── task/                # Async task contracts (Asynq)
│   │
│   ├── worker/              # Worker orchestration (processors)
│   │
│   ├── model/               # Domain models
│   │
│   ├── platform/            # Infrastructure (DB, Redis, config)
│   │
│   └── web/                 # Frontend hosting & dev proxy (optional)
│
├── frontend/                # Vite + React application
│
├── ARCHITECTURE.md
├── DETACH.md
└── README.md
```

---

## Backend Architecture

### Runtime Separation

The backend consists of **multiple runtimes**, each with a single responsibility:

* **Server runtime** (`cmd/server`)

  * Serves HTTP APIs
  * Optionally hosts the frontend
  * Enqueues async tasks

* **Worker runtime** (`cmd/worker`)

  * Consumes async tasks
  * Executes background jobs

Each runtime wires dependencies differently but shares the same core logic.

---

## Dependency Rules (Critical)

All dependencies flow inward toward the service layer.

```
handler → service → repository
worker  → service → repository
```

Forbidden dependencies:

* service → handler
* service → asynq
* service → echo
* repository → service
* frontend → backend internals

Violating these rules breaks testability and detachment.

---

## API Layer (`internal/api`)

**Purpose**: HTTP transport only.

Responsibilities:

* Parse HTTP requests
* Validate input
* Call services
* Serialize responses
* Enqueue async tasks

Non-responsibilities:

* Business rules
* Database logic
* Retry logic

Handlers should be thin and boring.

---

## Service Layer (`internal/service`)

**Purpose**: Core business logic.

Characteristics:

* Transport-agnostic
* Synchronous
* Deterministic
* Fully testable

Services:

* Accept primitives or domain objects
* Enforce business rules
* Call repositories
* Return domain errors

Services must not:

* Import Asynq
* Import HTTP frameworks
* Spawn goroutines

This layer is the heart of the system.

---

## Repository Layer (`internal/repository`)

**Purpose**: Data access abstraction.

Structure:

* Interfaces define required behavior
* Implementations live in subpackages (e.g. postgres)

Repositories:

* Perform persistence only
* Do not contain business logic
* Are shared by server and worker

---

## Asynchronous Architecture (Asynq)

### Task Contracts (`internal/task`)

Task definitions are shared contracts between server and worker.

They contain:

* Task type names
* Payload structures

They must not contain:

* Redis logic
* Asynq client/server code

This makes async transport replaceable.

---

### Server: Enqueueing Tasks

The server runtime:

* Creates tasks
* Enqueues them via Asynq client
* Never executes async work inline

This keeps request latency low and predictable.

---

### Worker: Processing Tasks

The worker runtime:

* Runs Asynq server
* Maps task types to processors
* Calls services

Processors:

* Deserialize payloads
* Call services
* Return errors for retry handling

Retry, backoff, and DLQ are handled by Asynq, not services.

---

## Platform Layer (`internal/platform`)

**Purpose**: Infrastructure wiring.

Contains:

* Database connections
* Redis connections
* Asynq client/server setup
* Configuration loading

This layer isolates external systems from business logic.

---

## Frontend Architecture

The frontend is a standalone Vite + React application.

### Directory Structure

```
frontend/
├── components/
│   ├── ui/          # shadcn components
│   └── layout/
│
├── hooks/            # data + behavior hooks
│
├── lib/              # API client, utilities, env
│
├── pages/            # route-level components
│
├── router/           # TanStack Router configuration
│
├── main.tsx
└── index.html
```

### Frontend Rules

* Communicates with backend only via HTTP (`/api`)
* No shared types with Go
* Routing handled by TanStack Router
* Data fetching via hooks
* UI primitives isolated in `components/ui`

---

## Frontend Hosting (Optional)

In development:

* Go server proxies requests to Vite dev server

In production:

* Frontend is embedded into Go binary

This hosting layer lives in `internal/web` and can be deleted without affecting APIs or workers.

---

## Detachment Guarantees

The architecture guarantees:

* Frontend can be deployed separately
* Worker can scale independently
* Async system can be replaced
* New runtimes can be added

Examples:

* Add gRPC: new `cmd/grpc`
* Add cron jobs: new `cmd/cron`
* Replace Asynq: modify `task/` + `worker/`

---

## Guiding Principles

* Clear ownership per layer
* Explicit runtime boundaries
* Transport-agnostic services
* Boring, predictable code

If a change forces you to break these rules, the architecture is being violated.
