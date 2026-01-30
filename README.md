# Go + Vite + React — Clean Full-Stack Example

This repository demonstrates a **clean, detachable integration** between:

* **Go (Echo)** backend API
* **Vite + React (TypeScript)** frontend
* Optional **single-binary deployment** via Go `embed`
* First-class **Vite HMR** during development

The key goal is to support a great developer experience **without coupling frontend and backend**, so they can be separated later with minimal changes.

---

## Why This Exists

I originally wanted to embed a Vite + React app into a Go binary to ship a single full-stack executable.

I started from ideas in [this great video](https://www.youtube.com/watch?v=w_Rv_3-FF0g), which shows how to embed frontend assets in Go. While that works, it compromises one of Vite’s biggest strengths: **its dev server and hot module reloading**.

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
├── internal/
│   ├── api/                # API routes & handlers
│   └── web/                # OPTIONAL web hosting layer
│       ├── web.go          # dev vs prod selection
│       ├── dev_proxy.go    # dev-only Vite proxy
│       ├── static.go       # static file serving
│       └── embedded_assets.go # go:embed (prod only)
│
├── frontend/               # Standalone Vite + React app
│   ├── src/
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

* `internal/api` → **application logic**
* `internal/web` → **optional infrastructure**

The backend does not depend on the frontend to function.

---

## Backend Overview (Go)

* Uses **Echo**
* Exposes APIs under `/api/*`
* Contains no frontend-specific logic
* Can be deployed independently as a pure API service

The backend may optionally:

* proxy frontend requests in development
* serve embedded static assets in production

These behaviors are **completely removable**.

---

## Frontend Overview (Vite + React)

The `frontend/` directory is a **normal Vite project**, created with:

```bash
yarn create vite
```

* Runs independently with `npm run dev`
* Uses Vite’s dev server and HMR
* Communicates with backend via HTTP only

### API Proxy (Frontend-Owned)

```ts
// vite.config.ts
server: {
  proxy: {
    '/api': 'http://localhost:3000'
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

* Go API server (optionally with live reload via `air`)
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
3. Produce a single executable

Run with:

```bash
./server
```

This mode is optional — the frontend can also be deployed separately.

---

## Docker Build

```bash
docker build -t go-vite .
```

Uses a multi-stage Dockerfile:

1. Node image to build frontend assets
2. Go image to compile backend with embedded assets
3. Final Alpine image containing a single binary

---

## Hot Module Reloading (How It Works)

### Development

* Frontend runs on Vite dev server (`localhost:5173`)
* Backend serves `/api/*`
* Optional proxy forwards non-API requests to Vite

This preserves **native Vite HMR behavior** exactly as intended.

### Production

* Frontend assets are embedded via `go:embed`
* Backend serves static files directly
* No Node or Vite runtime required

---

## Detaching Frontend and Backend Later

This repo is intentionally structured so detaching is trivial.

When ready:

1. Remove `internal/web`
2. Remove the web registration call in `main.go`
3. Enable CORS
4. Deploy frontend to a CDN
5. Deploy backend as an API service

No API handlers change.
No business logic changes.

See **`DETACH.md`** for a step-by-step guide.

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
