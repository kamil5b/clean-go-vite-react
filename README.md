# Go + Vite + React — Clean Full-Stack Example (Lite)

> **Note**: This is a **lite/simplified version** of the main clean-go-vite-react project, focusing on essential features and a streamlined architecture for rapid prototyping and small-to-medium applications.

This repository demonstrates a **clean, detachable integration** between:

* **Go (Echo)** backend API with simplified Clean Architecture
* **Vite + React (TypeScript)** frontend
* Optional **single-binary deployment** via Go `embed`
* First-class **Vite HMR** during development

The key goal is to support a great developer experience **without coupling frontend and backend**, so they can be separated later with minimal changes.

---

## Why This Exists

I originally wanted to embed a Vite + React app into a Go binary to ship a single full-stack executable.

I started from ideas in [this great video](https://www.youtube.com/watch?v=w_Rv_3-FF0g), which shows how to embed frontend assets in Go. While that works, it compromises one of Vite's biggest strengths: **its dev server and hot module reloading**.

This project keeps **Vite HMR exactly as-is in development**, while still allowing:

* embedded assets in production
* a single deployable binary (optional)
* clean separation between frontend and backend

---

## Lite Version Differences

This lite version provides a **simplified architecture** compared to the full version:

* **Streamlined layers**: `domain` (business logic + models), `infra` (infrastructure), and `api` (HTTP layer)
* **Faster setup**: Fewer abstractions, quicker to understand and modify
* **Production-ready foundations**: Still maintains clean separation and testability
* **Easy to scale up**: Can evolve into the full Clean Architecture when needed

Perfect for:
- MVPs and prototypes
- Small-to-medium applications
- Learning Go + React integration
- Projects that may need to scale architecture later

---

## Design Principles

This repo follows a few non-negotiable rules:

1. **The backend does not know React exists**
2. **The frontend does not know how the backend is implemented**
3. **Frontend integration is optional infrastructure, not core logic**

If you delete the frontend tomorrow, the backend still builds and runs.

---

## Project Structure

```txt
.
├── cmd/
│   └── server/
│       └── main.go          # Application entrypoint (composition only)
│
├── backend/                # Backend (Simplified Clean Architecture)
│   ├── api/               # HTTP layer (handlers, middleware, routes)
│   │   ├── handler/       # Request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── router.go      # Route definitions
│   ├── domain/            # Business logic + models
│   │   ├── logic.go       # Business logic implementation
│   │   └── models.go      # Data models/entities
│   └── infra/             # Infrastructure (database, config)
│       └── db.go          # Database initialization
│
├── embedder/              # OPTIONAL web hosting layer
│   └── embedder.go        # dev proxy + static serving
│
├── frontend/              # Standalone Vite + React app
│   ├── src/
│   │   ├── api/          # API client modules
│   │   ├── components/   # React components (UI library)
│   │   ├── contexts/     # React contexts (auth, etc)
│   │   ├── hooks/        # Custom React hooks
│   │   ├── pages/        # Page components
│   │   ├── router/       # Routing configuration
│   │   └── ...
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
│
├── scripts/               # Development and build scripts
├── Makefile
├── go.mod
├── main.go               # Legacy entrypoint (delegates to cmd/server)
└── README.md
```

### Important Distinction

* `backend/` → **application logic** (simplified 3-layer architecture)
* `embedder/` → **optional infrastructure**

The backend does not depend on the frontend to function.

---

## Backend Overview (Go)

* Uses **Echo** web framework
* Implements a **simplified Clean Architecture** with clear layer separation:
  - **domain**: Business logic and data models
  - **infra**: Infrastructure concerns (database, config)
  - **api**: HTTP layer (handlers, routes, middleware)
* Exposes APIs under `/api/*`
* Contains no frontend-specific logic
* Can be deployed independently as a pure API service

The backend may optionally:

* proxy frontend requests in development
* serve embedded static assets in production

These behaviors are **completely removable**.

---

## Features

### Authentication

* User registration and login with JWT tokens
* HTTP-only cookies for token storage
* Protected routes with authentication middleware
* Session management with `useAuth` hook

### Items CRUD

A complete Items management system demonstrating full-stack CRUD operations:

#### Backend API Endpoints

All item endpoints require authentication via JWT token.

- `POST /api/items` — Create a new item
  - Request: `{ "title": string, "description": string }`
  - Response: `ItemInfo`

- `GET /api/items` — List all items for the authenticated user
  - Response: `ItemInfo[]`

- `GET /api/items/:id` — Get a single item
  - Response: `ItemInfo`

- `PUT /api/items/:id` — Update an item
  - Request: `{ "title": string, "description": string }`
  - Response: `ItemInfo`

- `DELETE /api/items/:id` — Delete an item
  - Response: `204 No Content`

#### Item Model

```go
type Item struct {
    ID          uuid.UUID
    Title       string
    Description string
    UserID      uuid.UUID  // Ownership validation
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### Frontend Implementation

- **API Client** (`frontend/src/api/items.ts`) — RESTful API methods for all CRUD operations
- **Custom Hook** (`frontend/src/hooks/useItems.ts`) — State management for items with loading and error handling
- **Page Component** (`frontend/src/pages/ItemsPage.tsx`) — Full UI for creating, viewing, editing, and deleting items
- **UI Components** — Uses the project's component library (Card, Table, Button, Field) for consistent styling

#### Features

- ✅ Create items with title and description
- ✅ View all items in a sortable table
- ✅ Edit items inline with form validation
- ✅ Delete items with confirmation
- ✅ Automatic owner validation (users can only access their items)
- ✅ Responsive design with Tailwind CSS
- ✅ Error handling and loading states

---

## Frontend Overview (Vite + React)

The `frontend/` directory is a **normal Vite project**, created with:

```bash
yarn create vite
```

* Runs independently with `npm run dev`
* Uses Vite's dev server and HMR
* Communicates with backend via HTTP only
* Includes TypeScript, modern tooling, and strict type organization

### API Proxy (Frontend-Owned)

```ts
// vite.config.ts
server: {
  proxy: {
    '/api': 'http://localhost:8080'
  }
}
```

This allows frontend and backend to run on separate ports with no backend awareness of Vite.

---

## Development Mode

```bash
make dev
```

Runs:

* Go API server (with live reload via `air`)
* Vite dev server with full HMR

Frontend changes (TSX, CSS, state) update instantly **without page reload**.

---

## Production Build (Unified Binary)

```bash
make build
```

This will:

1. Build frontend assets into `frontend/dist`
2. Embed those assets into the Go binary
3. Produce a single executable in `bin/server`

Run with:

```bash
./bin/server
```

This mode is optional — the frontend can also be deployed separately.

---

## Hot Module Reloading (How It Works)

### Development

* Frontend runs on Vite dev server (`localhost:5173`)
* Backend serves `/api/*` on (`localhost:8080`)
* Optional proxy forwards non-API requests to Vite

This preserves **native Vite HMR behavior** exactly as intended.

### Production

* Frontend assets are embedded via `go:embed`
* Backend serves static files directly
* No Node or Vite runtime required

---

## Embedded Deployment with Easy Detachment

The Vite React app is **embedded into the Go binary** for single-file deployment, but the architecture makes it **trivially easy to detach later** when you need separate services.

**Key Benefit**: Start simple with one binary, scale to microservices without refactoring.

### How Embedding Works

In production, `frontend/dist` is embedded into the Go binary. The `embedder/` package handles:
- **Development**: Proxies requests to Vite dev server (preserves HMR)
- **Production**: Serves embedded static files from memory

This gives you a **single deployable binary** with the full stack.

---

### When to Detach

Detach when you need:
- Independent scaling of frontend and backend
- CDN distribution for frontend assets
- Separate deployment pipelines
- Multiple frontends consuming the same API

### How to Detach

Because the architecture is designed for detachment, it takes only a few steps:

#### Backend Changes

1. **Remove embedder integration** in `cmd/server/main.go`:
   ```go
   // Remove this line:
   e.Any("/*", echo.WrapHandler(web.Handler()))
   ```

2. **Enable proper CORS** in `cmd/server/main.go`:
   ```go
   e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
       AllowOrigins: []string{"https://your-frontend-domain.com"},
       AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
       AllowHeaders: []string{"Authorization", "Content-Type"},
       AllowCredentials: true,
   }))
   ```

3. **Delete embedder directory** (optional):
   ```bash
   rm -rf embedder/
   ```

4. **Deploy backend** as standalone API service

### Frontend Changes

1. **Update API base URL** in `frontend/vite.config.ts`:
   ```ts
   // Remove proxy in production, or configure for your API domain
   export default defineConfig({
     // ... other config
     server: {
       proxy: {
         '/api': 'http://localhost:8080' // Only for local dev
       }
     }
   })
   ```

2. **Set production API URL** in your frontend code:
   ```ts
   const API_BASE_URL = import.meta.env.PROD 
     ? 'https://your-backend-api.com'
     : '';
   ```

3. **Deploy frontend** to CDN (Vercel, Netlify, Cloudflare Pages, etc.)

#### Result

- ✅ **Backend**: Pure API service (e.g., `api.yourdomain.com`)
- ✅ **Frontend**: Static site on CDN (e.g., `yourdomain.com`)
- ✅ **Zero refactoring**: No changes to API handlers or business logic
- ✅ **Clean separation**: Frontend and backend communicate via HTTPS with CORS

**The entire detachment takes less than 10 minutes** because the architecture never couples them in the first place.

---

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Make (optional, but recommended)

### Quick Start

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd clean-go-vite-react
   ```

2. **Install frontend dependencies**:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

3. **Run development mode**:
   ```bash
   make dev
   ```

4. **Visit**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080/api

5. **Create an account** and start managing items!

---

## Available Commands

```bash
make dev          # Start development servers (Go + Vite)
make build        # Build production binary with embedded frontend
make build-fe     # Build frontend only
make build-be     # Build backend only
make clean        # Clean build artifacts
make test         # Run tests
```

---

## Notes & Gotchas

* Even in dev mode, the Go build requires embedded assets to exist.
  Run `make build-fe` once initially to populate `frontend/dist`.
* All API routes should live under `/api/*`.
* Avoid importing frontend artifacts into backend code.
* The `main.go` in root is deprecated - use `cmd/server/main.go` instead.
* All protected endpoints require authentication via JWT token in HTTP-only cookies or Authorization header.

---

## Scaling to Full Clean Architecture

When your application grows, you can evolve this lite version to the full Clean Architecture:

1. Split `domain/logic.go` into separate service layer
2. Extract repository interfaces and implementations from `domain`
3. Add dependency injection container
4. Introduce use cases / interactors
5. Add comprehensive testing layers

The simplified structure makes this evolution straightforward without major rewrites.

---

## Credits

Original inspiration from:

* [Embedding Vite into Go (YouTube)](https://www.youtube.com/watch?v=w_Rv_3-FF0g)
* Original example by @danhawkins

This lite fork focuses on **simplicity, speed, and clean boundaries** while maintaining the flexibility to scale.