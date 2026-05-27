# Phase 3: Security Priority Design

> **Goal:** Add refresh token rotation, login rate limiting, and CSP security headers to harden the application against brute-force, XSS, and token theft attacks.

**Architecture:** Short-lived access tokens (15min) + long-lived refresh tokens (7d, httpOnly cookie). IP-based rate limiting on login endpoint. Security headers middleware for all responses.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript (frontend), PostgreSQL

---

## 1. Refresh Token

### Current State
- Single JWT with 24h expiry
- Stored in localStorage (XSS vulnerable)
- Logout is client-side only (token remains valid)

### Design

| Token | Expiry | Storage | Usage |
|-------|--------|---------|-------|
| Access Token | 15 minutes | localStorage | API authentication |
| Refresh Token | 7 days | httpOnly cookie | Token refresh |

### Backend Changes

**`internal/library/jwt/jwt.go`:**
- Add `GenerateRefreshToken(userId uint64) (string, error)` method
- Refresh token only contains `user_id` (no roles/permissions — lightweight)
- Different expiry (7 * 24h)

**`internal/controller/auth.go`:**
- `Login`: Generate both tokens, set refresh_token as httpOnly cookie
- `Refresh`: Validate refresh_token from cookie, issue new access_token + refresh_token
- `Logout`: Clear refresh_token cookie

**`api/v1/auth.go`:**
- Update `LoginRes` to include `access_token` (rename from `token`)
- Add `RefreshReq` (empty, reads from cookie)
- Add `RefreshRes` with new `access_token`
- Add `LogoutRes` (empty)

**`internal/cmd/cmd.go`:**
- Add `POST /api/v1/refresh` route (no Auth middleware — reads cookie)
- Add `POST /api/v1/logout` route (requires Auth middleware)

### Frontend Changes

**`web/src/utils/request.ts`:**
- Add 401 interceptor: call `/api/v1/refresh`, retry original request
- If refresh fails (401), redirect to login

**`web/src/api/auth.ts`:**
- Add `refresh()` method
- Update `LoginRes` to use `access_token` field

**`web/src/store/modules/user.ts`:**
- Update `setToken` to use `access_token` field from response

### Cookie Settings

```
refresh_token: <token_value>
HttpOnly: true
Secure: false (dev) / true (prod)
SameSite: Lax
Path: /api/v1
Max-Age: 604800 (7 days)
```

---

## 2. Login Rate Limiting

### Current State
- No rate limiting on any endpoint
- Brute-force attacks are trivial

### Design

**In-memory rate limiter** (no Redis dependency):

```go
type RateLimiter struct {
    mu       sync.Mutex
    attempts map[string][]time.Time
    maxCount int
    window   time.Duration
}

func NewRateLimiter(maxCount int, window time.Duration) *RateLimiter
func (rl *RateLimiter) Allow(key string) bool
```

**Configuration:**
| Parameter | Value |
|-----------|-------|
| Window | 1 minute |
| Max attempts | 5 per IP |
| Exceeded response | HTTP 429 |

**Application:** Only on `POST /api/v1/login` route.

**Key:** Client IP address (from `r.GetClientIp()`).

---

## 3. CSP Security Headers

### Current State
- Only CORS headers (`Access-Control-Allow-Origin: *`)
- No security headers

### Design

New middleware `SecurityHeaders` applied globally:

```go
func SecurityHeaders(r *ghttp.Request) {
    r.Response.Header().Set("X-Content-Type-Options", "nosniff")
    r.Response.Header().Set("X-Frame-Options", "DENY")
    r.Response.Header().Set("X-XSS-Protection", "1; mode=block")
    r.Response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
    r.Response.Header().Set("Content-Security-Policy", 
        "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
    r.Middleware.Next()
}
```

**Headers:**
| Header | Value | Purpose |
|--------|-------|---------|
| X-Content-Type-Options | nosniff | Prevent MIME sniffing |
| X-Frame-Options | DENY | Prevent clickjacking |
| X-XSS-Protection | 1; mode=block | Enable browser XSS filter |
| Referrer-Policy | strict-origin-when-cross-origin | Control referrer info |
| Content-Security-Policy | (see above) | Restrict resource sources |

**Application:** Global middleware (same level as CORS).

---

## File Changes Summary

### New Files
| File | Purpose |
|------|---------|
| `internal/middleware/ratelimit.go` | Rate limiter middleware |
| `internal/middleware/security.go` | Security headers middleware |

### Modified Files
| File | Change |
|------|--------|
| `internal/library/jwt/jwt.go` | Add GenerateRefreshToken |
| `internal/controller/auth.go` | Login returns 2 tokens; add Refresh/Logout |
| `api/v1/auth.go` | Add RefreshReq/Res, update LoginRes |
| `internal/cmd/cmd.go` | Add /refresh route, SecurityHeaders middleware, rate limiter on login |
| `web/src/utils/request.ts` | 401 interceptor with refresh |
| `web/src/api/auth.ts` | Add refresh(), update LoginRes |
| `web/src/store/modules/user.ts` | Update token field name |

---

## Migration

No database migration needed — refresh tokens are stateless JWTs validated by signature and expiry.

---

## Security Improvements Summary

| Before | After |
|--------|-------|
| Single 24h token | 15min access + 7d refresh |
| Token in localStorage only | Access in localStorage, refresh in httpOnly cookie |
| No token rotation | Refresh rotates both tokens |
| No brute-force protection | 5 attempts/minute per IP |
| No security headers | 5 security headers (CSP, X-Frame, etc.) |
| Logout is client-only | Server-side cookie clearing |
