# Authentication System

This document describes the authentication system for the Go + Vite + React application.

## Overview

The authentication system is designed as a **pure API concern** with the frontend as a dumb HTTP client. It uses:

- **HTTP-only cookies** for token storage (no localStorage, no XSS vectors)
- **JWT tokens** for stateless authentication
- **CSRF protection** with `SameSite=Lax` and token-based validation
- **Bcrypt password hashing** with default cost factor

This approach works identically in:
- Vite dev mode with proxy
- Embedded single-binary production
- Detached frontend + API later

## Architecture Principles

### Design Goals

1. **Backend owns auth entirely** — All logic resides in the API
2. **Frontend is a dumb HTTP client** — Uses standard `fetch` with `credentials: 'include'`
3. **No special dev-only hacks** — Same auth flow everywhere
4. **Easy to detach later** — Zero coupling to frontend framework

### Why HTTP-Only Cookies?

| Concern | Cookies | Bearer Tokens |
|---------|---------|---------------|
| Vite dev proxy | ✅ Automatic | ⚠️ CORS headers needed |
| Embedded prod | ✅ Seamless | ✅ Works |
| XSS resistance | ✅ HttpOnly flag | ❌ Token readable |
| CSRF protection | ✅ Built-in | ⚠️ Manual handling |
| Detaching frontend | ✅ No changes | ⚠️ More config |
| Mobile/CLI clients | ⚠️ Less ideal | ✅ Better |

For this **web-first** repo, cookies are the pragmatic choice.

## API Endpoints

### Public Endpoints

#### `POST /api/auth/register`

Register a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secure-password",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cookies Set:**
- `access_token` (JWT, 15 minutes)
- `refresh_token` (JWT, 7 days)

**Error Responses:**
- `400 Bad Request` — Missing or invalid fields
- `409 Conflict` — Email already registered

---

#### `POST /api/auth/login`

Authenticate an existing user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secure-password"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cookies Set:**
- `access_token` (JWT, 15 minutes)
- `refresh_token` (JWT, 7 days)

**Error Responses:**
- `400 Bad Request` — Missing email or password
- `401 Unauthorized` — Invalid credentials

---

#### `POST /api/auth/refresh`

Generate a new access token using the refresh token.

**Request:**
No body required. Refresh token is sent via cookie.

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cookies Set:**
- `access_token` (JWT, 15 minutes, replaces old)

**Error Responses:**
- `401 Unauthorized` — Missing or invalid refresh token

---

#### `GET /api/csrf`

Get a CSRF token for state-changing operations.

**Request:**
No parameters required.

**Response (200 OK):**
```json
{
  "token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

---

### Protected Endpoints

All protected endpoints require a valid `access_token` cookie.

#### `GET /api/auth/me`

Get the current authenticated user's information.

**Request Headers:**
- Cookie: `access_token=...` (automatic via browser)

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Error Responses:**
- `401 Unauthorized` — Missing or invalid token

---

#### `POST /api/auth/logout`

Clear authentication cookies and logout the user.

**Request Headers:**
- Cookie: `access_token=...` (automatic via browser)

**Response (200 OK):**
```json
{
  "message": "logged out successfully"
}
```

**Cookies Cleared:**
- `access_token` (MaxAge: -1)
- `refresh_token` (MaxAge: -1)

---

## Token Details

### Access Token (JWT)

**Type:** JWT (HMAC-SHA256)  
**Expiry:** 15 minutes  
**Storage:** HTTP-only cookie  
**Purpose:** Authenticate API requests

**Claims:**
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "iat": 1699500000,
  "exp": 1699500900,
  "iss": "go-vite-react"
}
```

### Refresh Token (JWT)

**Type:** JWT (HMAC-SHA256)  
**Expiry:** 7 days  
**Storage:** HTTP-only cookie  
**Purpose:** Obtain new access tokens

**Claims:**
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "iat": 1699500000,
  "exp": 1699608000,
  "iss": "go-vite-react"
}
```

---

## Cookies

### Access Token Cookie

```http
Set-Cookie: access_token=eyJ...; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=900
```

- **HttpOnly:** Prevents JavaScript access (XSS protection)
- **Secure:** Only sent over HTTPS (set in production)
- **SameSite=Lax:** CSRF protection for same-site requests
- **Max-Age=900:** 15 minutes

### Refresh Token Cookie

```http
Set-Cookie: refresh_token=eyJ...; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=604800
```

- **HttpOnly:** Prevents JavaScript access
- **Secure:** Only sent over HTTPS (set in production)
- **SameSite=Lax:** CSRF protection
- **Max-Age=604800:** 7 days

---

## CSRF Protection

### Strategy

**SameSite=Lax** covers most cases. For additional protection on sensitive operations:

1. **Frontend requests CSRF token:**
   ```typescript
   const csrf = await fetch('/api/csrf').then(r => r.json());
   ```

2. **Frontend sends token in header:**
   ```typescript
   fetch('/api/users/update', {
     method: 'POST',
     credentials: 'include',
     headers: {
       'X-CSRF-Token': csrf.token
     },
     body: JSON.stringify({ ... })
   })
   ```

3. **Backend validates token:**
   - Middleware checks `X-CSRF-Token` header
   - Validates token (currently basic validation; extend as needed)

### When CSRF Headers Are Required

State-changing operations require `X-CSRF-Token` header:
- `POST` requests
- `PUT` requests
- `PATCH` requests
- `DELETE` requests

Safe operations do NOT require CSRF tokens:
- `GET` requests
- `HEAD` requests
- `OPTIONS` requests

---

## Environment Configuration

### Required Environment Variables

Set these in `.env` or your deployment platform:

```bash
# JWT Secrets (MUST change in production!)
JWT_ACCESS_SECRET="your-access-token-secret-key-change-in-prod"
JWT_REFRESH_SECRET="your-refresh-token-secret-key-change-in-prod"
```

### Recommended Environment Variables

```bash
# Database
DATABASE_DSN="dev.db"
DATABASE_TYPE="sqlite"

# Server
SERVER_PORT=8080
SERVER_HOST=
```

### Development vs Production

**Development (.env):**
```bash
JWT_ACCESS_SECRET="dev-access-secret"
JWT_REFRESH_SECRET="dev-refresh-secret"
```

**Production:**
```bash
# Use strong, randomly generated secrets
# Example: openssl rand -base64 32

JWT_ACCESS_SECRET="$(openssl rand -base64 32)"
JWT_REFRESH_SECRET="$(openssl rand -base64 32)"
```

For production with HTTPS, also set in `backend/api/middleware/auth.go`:
```go
// Change Secure: false to Secure: true
c.SetCookie(&http.Cookie{
  Secure: true,  // Only send over HTTPS
})
```

---

## Frontend Integration

### Basic Setup

Frontend uses standard `fetch` API with `credentials: 'include'` to send cookies automatically.

### Registration

```typescript
async function register(email: string, password: string, name: string) {
  const response = await fetch('/api/auth/register', {
    method: 'POST',
    credentials: 'include',  // Send cookies
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, name })
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return response.json();
}
```

### Login

```typescript
async function login(email: string, password: string) {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    credentials: 'include',  // Send cookies
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return response.json();
}
```

### Get Current User

```typescript
async function getCurrentUser() {
  const response = await fetch('/api/auth/me', {
    credentials: 'include'  // Send cookies
  });

  if (!response.ok) {
    if (response.status === 401) {
      return null;  // Not authenticated
    }
    throw new Error('Failed to fetch user');
  }

  return response.json();
}
```

### Refresh Token

```typescript
async function refreshToken() {
  const response = await fetch('/api/auth/refresh', {
    method: 'POST',
    credentials: 'include'  // Send refresh_token cookie
  });

  if (!response.ok) {
    // Refresh failed, user needs to login again
    return null;
  }

  return response.json();
}
```

### Logout

```typescript
async function logout() {
  await fetch('/api/auth/logout', {
    method: 'POST',
    credentials: 'include'
  });

  // Redirect to login page
  window.location.href = '/login';
}
```

### CSRF Protection for State-Changing Requests

```typescript
async function makeStateChangingRequest(method: 'POST' | 'PUT' | 'DELETE', url: string, data?: any) {
  // Get CSRF token
  const csrfResponse = await fetch('/api/csrf');
  const { token: csrfToken } = await csrfResponse.json();

  const response = await fetch(url, {
    method,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken
    },
    body: data ? JSON.stringify(data) : undefined
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return response.json();
}
```

---

## Development Workflow

### Starting the App

```bash
# Terminal 1: Backend
go run ./cmd/server

# Terminal 2: Frontend (in another terminal)
cd frontend
npm run dev
```

The auth flow works exactly the same as production:
- Vite proxy forwards `/api/*` to Go backend
- Cookies are handled automatically by the browser
- No special dev configuration needed

### Testing with cURL

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","name":"Test User"}' \
  -c cookies.txt

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}' \
  -c cookies.txt

# Get current user (using saved cookies)
curl -X GET http://localhost:8080/api/auth/me \
  -b cookies.txt

# Logout
curl -X POST http://localhost:8080/api/auth/logout \
  -b cookies.txt
```

---

## Security Considerations

### Password Security

- Passwords are hashed with **bcrypt** (default cost: 10)
- Never stored in plaintext
- Validated on every login attempt

### Token Security

- Tokens are **signed with HMAC-SHA256**
- Cannot be modified without the secret
- **Access tokens expire after 15 minutes**
- **Refresh tokens expire after 7 days**

### Cookie Security

- **HttpOnly flag prevents XSS attacks** — JavaScript cannot access cookies
- **SameSite=Lax prevents CSRF attacks** — Cookies not sent cross-site
- **Secure flag (production only)** — Cookies only sent over HTTPS

### Additional Protections

- CSRF token validation for state-changing requests
- Context cancellation awareness (respects request timeouts)
- Clear error messages without leaking sensitive info

### Production Checklist

- [ ] **Change JWT secrets** — Generate new values with `openssl rand -base64 32`
- [ ] **Enable HTTPS** — Set `Secure: true` on cookies
- [ ] **Set strong secrets** — At least 32 bytes of randomness
- [ ] **Enable rate limiting** — On `/api/auth/login` and `/api/auth/register`
- [ ] **Monitor failed logins** — Detect brute force attempts
- [ ] **Use HTTPS in Vite** — For prod-like testing
- [ ] **Configure CORS** — If frontend is detached

---

## Detaching Frontend and Backend Later

The auth system is already designed for easy detachment:

### Step 1: Deploy Frontend Separately

No changes to auth code needed. The frontend can be deployed to any host.

### Step 2: Configure CORS

In `backend/api/router.go`, add CORS middleware:

```go
import "github.com/labstack/echo/v4/middleware"

func SetupRoutes(...) {
  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://yourdomain.com"},
    AllowCredentials: true,
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
  }))
  // ... rest of routes
}
```

### Step 3: Update Frontend API Base URL

```typescript
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

fetch(`${API_BASE_URL}/api/auth/login`, { ... })
```

**That's it.** No auth logic changes needed. The API works the same way.

---

## Extending the System

### Storing Refresh Tokens for Revocation

Currently, refresh tokens are validated via JWT claims only. For revocation support (logout, device removal, etc.):

1. Create a `RefreshTokenEntity` with database storage (already defined in `backend/model/entity/refresh_token.go`)
2. Implement `RefreshTokenRepository` interface
3. Store issued refresh tokens in the database
4. Check token validity on `/api/auth/refresh` and `/api/auth/logout`

### Rate Limiting

Add to login/register endpoints:

```go
import "github.com/labstack/echo-contrib/echoprometheus"

e.Use(middleware.RateLimiter(...))
```

### Multi-Device Sessions

Track active sessions per user:

```go
type SessionEntity struct {
  ID        uuid.UUID
  UserID    uuid.UUID
  Token     string
  Device    string
  ExpiresAt time.Time
}
```

### Two-Factor Authentication

After login succeeds, challenge the user:
- `POST /api/auth/challenge/2fa` — Send 2FA code
- `POST /api/auth/verify/2fa` — Verify and issue tokens

### OAuth/Social Login

Add social auth providers:
- `POST /api/auth/github` — Redirect to GitHub OAuth
- `POST /api/auth/callback` — Handle OAuth callback

All while keeping the same HTTP-only cookie response format.

---

## Troubleshooting

### "missing access token" on protected endpoints

**Problem:** Frontend not sending cookies with requests.

**Solution:** Ensure `credentials: 'include'` in fetch calls:
```typescript
fetch('/api/auth/me', {
  credentials: 'include'  // This is required
})
```

### "invalid email or password" even with correct credentials

**Problem:** Wrong hashing or comparison.

**Solutions:**
- Check bcrypt cost factor matches (default: 10)
- Verify password field is not trimmed unexpectedly
- Ensure database stores bcrypt hash correctly

### Cookies not being set in production

**Problem:** `Secure` flag set but not using HTTPS.

**Solutions:**
- Enable HTTPS in production
- During testing, set `Secure: false` (dev only)
- Use proper SSL certificates

### CSRF token errors on state-changing requests

**Problem:** Missing or invalid `X-CSRF-Token` header.

**Solutions:**
- Fetch CSRF token first: `GET /api/csrf`
- Send token in header: `'X-CSRF-Token': token`
- Ensure header name matches exactly (case-sensitive)

### Token expired but app still shows logged in

**Problem:** Frontend not calling refresh endpoint.

**Solution:** Implement token refresh logic:
```typescript
if (response.status === 401) {
  // Token expired, try refresh
  await refreshToken();
  // Retry original request
}
```

---

## References

- **OWASP Authentication Cheat Sheet:** https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
- **JWT Best Practices:** https://tools.ietf.org/html/rfc8949
- **HTTP Cookie Security:** https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies

---

## Support

For questions or issues with authentication:

1. Check the troubleshooting section above
2. Review example frontend integration code
3. Enable debug logging in token validation
4. Check JWT claims with `jwt.io` (for development only)