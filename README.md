# Go + Vite + React — Clean Full-Stack Example

A single-binary full-stack app that keeps frontend and backend cleanly separated, with an explicit and easy detachment path.

This repository demonstrates a **clean, detachable integration** between:

* **Go** backend API with Clean Architecture
* **Vite + React (TypeScript)** frontend
* **Single-binary deployment** via Go `embed` that easy to detach
* First-class **Vite HMR** during development
* **Framework-agnostic** frontend and backend serving layer

The key goal is to support a great developer experience **without coupling frontend and backend**, so they can be separated later with minimal changes.

<img width="5197" height="6315" alt="Project System" src="https://github.com/user-attachments/assets/9099ce9e-be10-481a-a006-b8711ca3457b" />

---

## Why This Exists

I originally wanted to embed a Vite + React app into a Go binary to ship a single full-stack executable.

I started from ideas in [this great video](https://www.youtube.com/watch?v=w_Rv_3-FF0g), which shows how to embed frontend assets in Go. While that works, it compromises one of Vite's biggest strengths: **its dev server and hot module reloading**.

This project keeps **Vite HMR exactly as-is in development**, while still allowing:

* embedded assets in production
* a single deployable binary (optional)
* clean separation between frontend and backend

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
├── backend/                # Backend Clean Architecture
│   ├── api/               # HTTP layer (handlers, middleware, routes)
│   ├── service/           # Business logic layer
│   ├── repository/        # Data access layer (interfaces + implementations)
│   ├── model/             # Data models (entities, DTOs)
│   ├── di/                # Dependency injection
│   └── platform/          # Infrastructure (config, database)
│
├── embedder/              # OPTIONAL framework-agnostic web hosting layer
│   └── embedder.go        # dev proxy + static serving (http.Handler)
│
├── frontend/              # Standalone Vite + React app
│   ├── src/
│   │   ├── api/          # API client modules
│   │   ├── types/        # TypeScript type definitions
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom React hooks
│   │   ├── pages/        # Page components
│   │   └── ...
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
│
├── Makefile
├── Dockerfile
├── go.mod
└── README.md
```

### Important Distinction

* `backend/` → **application logic**
* `embedder/` → **optional infrastructure**

The backend does not depend on the frontend to function.

---

## Backend Overview (Go)

* **Framework-agnostic** — use any router (Echo, Gin, Chi, standard `net/http`)
* Implements **Clean Architecture** with strict layer separation
* Exposes APIs under `/api/*`
* Contains no frontend-specific logic
* Can be deployed independently as a pure API service
* Includes JWT authentication, CSRF protection, and comprehensive testing
* **Full CRUD implementation** for Items, Tags, and Invoices with relationships
* **UUID-based primary keys** for all entities (except auto-increment for Items, Tags, and Invoices)

The backend may optionally use the `embedder` package to:

* proxy frontend requests in development
* serve embedded static assets in production

The embedder uses standard `http.Handler` interface and works with any Go web framework.

These behaviors are **completely removable**.

**For detailed backend documentation, see [`backend/README.md`](./backend/README.md)**

---

## Frontend Overview (Vite + React)

The `frontend/` directory is a **Vite project with Tailwind support**, created with:

```bash
yarn create vite
```

* Runs independently with `npm run dev`
* Uses Vite's dev server and HMR
* Communicates with backend via HTTP only
* Includes TypeScript, Tailwind (we use Shadcn), Zod validation, and strict type organization

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

### Key Features

* **Multi-select dropdowns** with infinite scroll for Items and Tags
* **Real-time search** with debounced API calls
* **Paginated tables** for all CRUD operations
* **Modal forms** for Item and Tag management
* **Full-page forms** for Invoice management with line items
* **Type-safe API client** with Zod validation
* Strict type organization: all API types defined in `src/types/`

**For detailed frontend documentation, see [`frontend/README.md`](./frontend/README.md)**

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

## Docker Build

```bash
docker build -t go-vite-react .
```

Uses a multi-stage Dockerfile:

1. Node image to build frontend assets
2. Go image to compile backend with embedded assets
3. Final Alpine image containing a single binary

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

In production, `frontend/dist` is embedded into the Go binary. The `embedder/` package is **framework-agnostic** and returns a standard `http.Handler` that:
- **Development**: Proxies requests to Vite dev server (preserves HMR)
- **Production**: Serves embedded static files from memory

This gives you a **single deployable binary** with the full stack.

**Integration Example** (works with any framework):
```go
// With Echo
e.Any("/*", echo.WrapHandler(embedder.Handler()))

// With standard net/http
mux.Handle("/", embedder.Handler())

// With Gin
r.NoRoute(gin.WrapH(embedder.Handler()))
```

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
   // Remove this line (example with Echo):
   e.Any("/*", echo.WrapHandler(embedder.Handler()))
   
   // Or with standard net/http:
   mux.Handle("/", embedder.Handler())
   ```

2. **Enable CORS** for your framework:
   ```go
   e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
       AllowOrigins: []string{"https://your-frontend-domain.com"},
       AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
       AllowHeaders: []string{"Authorization", "Content-Type", "X-CSRF-Token"},
       AllowCredentials: true,
   }))
   ```

3. **Delete embedder directory**:
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

## Features

### Implemented CRUD Operations

The application includes complete CRUD functionality for:

#### Items
- Create, Read, Update, Delete operations
- Pagination with configurable page size
- Search by name
- Modal-based forms
- Soft deletes with `deleted_at` timestamp

#### Tags
- Create, Read, Update, Delete operations
- Pagination with configurable page size
- Search by name
- Color picker with hex code validation
- Modal-based forms
- Used for categorizing invoices

#### Invoices
- Create, Read, Update, Delete operations
- Pagination with configurable page size
- Search by ID
- **Line items** with quantity, unit price, and total calculation
- **Many-to-many relationship** with Tags
- Multi-select dropdowns with infinite scroll for selecting items and tags
- Automatic grand total calculation
- Full-page forms with item table editor
- Detailed view page showing all invoice information

### Data Models

All entities use **UUID for user-related data** (authentication) and **auto-increment integers** for business entities (Items, Tags, Invoices). Key relationships:

- **Invoice → Invoice Items** (one-to-many)
- **Invoice → Tags** (many-to-many via junction table)
- **Invoice Item → Item** (many-to-one)

### API Endpoints

All endpoints are protected with JWT authentication and CSRF protection on mutations:

```
# Items
GET    /api/items              # List with pagination & search
POST   /api/items              # Create (CSRF protected)
GET    /api/items/:id          # Get by ID
PUT    /api/items/:id          # Update (CSRF protected)
DELETE /api/items/:id          # Delete (CSRF protected)

# Tags
GET    /api/tags               # List with pagination & search
POST   /api/tags               # Create (CSRF protected)
GET    /api/tags/:id           # Get by ID
PUT    /api/tags/:id           # Update (CSRF protected)
DELETE /api/tags/:id           # Delete (CSRF protected)

# Invoices
GET    /api/invoices           # List with pagination & search
POST   /api/invoices           # Create with items & tags (CSRF protected)
GET    /api/invoices/:id       # Get with all relations
PUT    /api/invoices/:id       # Update (replaces items & tags) (CSRF protected)
DELETE /api/invoices/:id       # Delete (CSRF protected)
```

## Documentation

### Other Guides
- **[AUTH.md](./AUTH.md)** - Authentication implementation guide

### Component Documentation
- **[backend/README.md](./backend/README.md)** - Complete backend architecture guide (Clean Architecture, TDD workflow)
- **[frontend/README.md](./frontend/README.md)** - Frontend structure and type rules

---

## Notes & Gotchas

* Even in dev mode, the Go build requires embedded assets to exist.
  Run `make build` once initially to populate `frontend/dist`.
* All API routes should live under `/api/*`.
* Avoid importing frontend artifacts into backend code.

---

## Credits

Original inspiration from:

* [Embedding Vite into Go (YouTube)](https://www.youtube.com/watch?v=w_Rv_3-FF0g)
* Original example by @danhawkins

This fork focuses on **clean boundaries and long-term maintainability**.
