# DETACH.md

## Purpose

This document describes **how to detach the frontend and backend** in this repository.

Detaching means:

* The **backend** becomes a standalone API service
* The **frontend** is deployed independently (CDN / static hosting)
* No business logic or API handlers are rewritten

If this process feels boring and mechanical, the architecture has done its job.

---

## Preconditions

Before detaching, ensure:

* All backend routes are under `/api/*`
* Frontend communicates with backend only via HTTP
* No HTML rendering or frontend logic exists in backend handlers

If these conditions are met, detaching is safe.

---

## High-Level Overview

Detaching consists of **four actions**:

1. Remove the optional web hosting layer from the backend
2. Enable CORS on the backend
3. Deploy the frontend separately
4. Point the frontend at the new API base URL

---

## Step 1 — Remove Web Hosting From Backend

### Delete Infrastructure Files

Remove the entire optional web layer:

```txt
internal/web/
```

This includes:

* Static asset serving
* Embedded frontend assets
* Dev proxy logic

These are no longer needed once the frontend is external.

---

### Update `main.go`

Before:

```go
web.Register(e, mode)
```

After:

```go
// frontend is now detached
```

No other backend code should change.

If this breaks compilation, it indicates frontend concerns leaked into backend logic and should be refactored.

---

## Step 2 — Enable CORS

Since the frontend is now hosted on a different origin, enable CORS.

Example (Echo):

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{
        "https://your-frontend-domain.com",
    },
    AllowMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPut,
        http.MethodDelete,
    },
}))
```

Do **not** use wildcard origins in production unless you fully understand the implications.

---

## Step 3 — Deploy Backend API

The backend can now be deployed independently:

* Containerized service
* VM / bare metal
* Managed platform

The backend binary no longer contains frontend assets and does not require Node or Vite.

Example API base URL:

```
https://api.example.com
```

---

## Step 4 — Deploy Frontend

### Build Frontend

From the `frontend/` directory:

```bash
npm run build
```

This produces static assets in `dist/`.

---

### Deploy Frontend

Deploy `dist/` to any static host:

* S3 + CloudFront
* Cloudflare Pages
* Vercel
* Netlify

The frontend is now fully independent.

---

### Update API Base URL

Configure the frontend to point at the new backend API.

Common approaches:

* Environment variables at build time
* Runtime configuration via `window.__ENV__`
* Reverse proxy at CDN layer

Example:

```ts
const API_BASE = import.meta.env.VITE_API_BASE_URL
```

---

## Validation Checklist

After detaching, verify:

* [ ] Backend starts without frontend code present
* [ ] Frontend loads from CDN
* [ ] API requests succeed cross-origin
* [ ] No `/` routes are handled by backend
* [ ] `/api/*` routes behave identically to before

If all checks pass, detachment is complete.

---

## Rollback Strategy

If detachment needs to be reversed:

1. Restore `internal/web`
2. Re-enable `web.Register(...)` in `main.go`
3. Rebuild unified binary

No data or API changes are required.

---

## Design Intent (Why This Is Easy)

This repo was structured so that:

* Frontend integration is **optional infrastructure**
* Backend logic never imports frontend artifacts
* Dev conveniences are removable without refactoring

Detachment should feel like **removing a feature flag**, not rewriting an application.

---

## Final Note

If detaching required touching API handlers, business logic, or domain models, stop and refactor before proceeding.

A clean boundary is the difference between an easy migration and a multi-week rewrite.
