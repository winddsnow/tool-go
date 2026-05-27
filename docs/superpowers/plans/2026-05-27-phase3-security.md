# Phase 3: Security Priority Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add refresh token rotation, login rate limiting, and CSP security headers.

**Architecture:** Short-lived access tokens (15min) + long-lived refresh tokens (7d, httpOnly cookie). IP-based rate limiting on login endpoint. Security headers middleware for all responses.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript (frontend)

---

## File Structure

### New Files
| File | Purpose |
|------|---------|
| `internal/middleware/ratelimit.go` | Rate limiter middleware |
| `internal/middleware/security.go` | Security headers middleware |

### Modified Files
| File | Change |
|------|--------|
| `internal/library/jwt/jwt.go` | Add GenerateRefreshToken method |
| `api/v1/auth.go` | Add RefreshReq/Res, update LoginRes field name |
| `internal/controller/auth.go` | Login returns 2 tokens; add Refresh/Logout handlers |
| `internal/cmd/cmd.go` | Add /refresh route, SecurityHeaders middleware, rate limiter on login |
| `web/src/utils/request.ts` | 401 interceptor with auto-refresh |
| `web/src/api/auth.ts` | Add refresh(), update LoginRes type |
| `web/src/store/modules/user.ts` | Update token field handling |

---

## Task 1: Backend — Security Headers Middleware

**Files:**
- Create: `internal/middleware/security.go`

- [ ] **Step 1: Create security headers middleware**

```go
package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// SecurityHeaders sets common security headers on all responses.
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

- [ ] **Step 2: Commit**

```bash
git add internal/middleware/security.go
git commit -m "feat: add security headers middleware (CSP, X-Frame, etc.)"
```

---

## Task 2: Backend — Rate Limiter Middleware

**Files:**
- Create: `internal/middleware/ratelimit.go`

- [ ] **Step 1: Create rate limiter**

```go
package middleware

import (
    "sync"
    "time"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

// RateLimiter provides IP-based rate limiting using an in-memory sliding window.
type RateLimiter struct {
    mu       sync.Mutex
    attempts map[string][]time.Time
    maxCount int
    window   time.Duration
}

// NewRateLimiter creates a rate limiter with maxCount requests per window.
func NewRateLimiter(maxCount int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        attempts: make(map[string][]time.Time),
        maxCount: maxCount,
        window:   window,
    }
}

// Allow checks if a request from the given key (IP) is allowed.
func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-rl.window)

    // Remove old entries
    attempts := rl.attempts[key]
    valid := make([]time.Time, 0, len(attempts))
    for _, t := range attempts {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }

    if len(valid) >= rl.maxCount {
        rl.attempts[key] = valid
        return false
    }

    rl.attempts[key] = append(valid, now)
    return true
}

// RateLimit returns a middleware that limits requests by client IP.
func RateLimit(maxCount int, window time.Duration) func(r *ghttp.Request) {
    limiter := NewRateLimiter(maxCount, window)
    return func(r *ghttp.Request) {
        ip := r.GetClientIp()
        if !limiter.Allow(ip) {
            g.Log().Warningf(r.GetCtx(), "rate limit exceeded for IP: %s", ip)
            r.Response.WriteJsonExit(g.Map{
                "code":    429,
                "message": "请求过于频繁，请稍后再试",
            })
            return
        }
        r.Middleware.Next()
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/middleware/ratelimit.go
git commit -m "feat: add IP-based rate limiter middleware for login endpoint"
```

---

## Task 3: Backend — Refresh Token in JWT

**Files:**
- Modify: `internal/library/jwt/jwt.go`

- [ ] **Step 1: Add refresh token claims and generation**

Read the current file first. Then add:

1. Add a new `RefreshClaims` struct:
```go
type RefreshClaims struct {
    jwt.RegisteredClaims
    UserId uint64 `json:"user_id"`
}
```

2. Add `GenerateRefreshToken` method:
```go
func (j *JWT) GenerateRefreshToken(userId uint64) (string, error) {
    claims := RefreshClaims{
        UserId: userId,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expires * 28)), // 28x access expiry = ~7 days
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secret))
}
```

3. Add `ParseRefreshToken` method:
```go
func (j *JWT) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.secret), nil
    })
    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorExpired != 0 {
                return nil, ErrTokenExpired
            }
            if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
                return nil, ErrTokenNotValidYet
            }
        }
        return nil, ErrTokenInvalid
    }
    if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, ErrTokenInvalid
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/library/jwt/jwt.go
git commit -m "feat: add refresh token generation and parsing to JWT library"
```

---

## Task 4: Backend — Update Auth API Types

**Files:**
- Modify: `api/v1/auth.go`

- [ ] **Step 1: Update LoginRes and add Refresh/Logout types**

Read the current file first.

1. Rename `Token` field in `LoginRes` to `AccessToken`:
```go
type LoginRes struct {
    AccessToken  string     `json:"access_token" dc:"访问令牌"`
    UserId       uint64     `json:"user_id" dc:"用户ID"`
    Username     string     `json:"username" dc:"用户名"`
    Nickname     string     `json:"nickname" dc:"昵称"`
    Roles        []string   `json:"roles" dc:"角色列表"`
    Menus        []MenuTree `json:"menus" dc:"菜单树"`
    Permissions  []string   `json:"permissions" dc:"权限码列表"`
}
```

2. Add RefreshReq and RefreshRes:
```go
type RefreshReq struct {
    g.Meta `path:"/refresh" method:"post" tags:"Auth" summary:"刷新访问令牌"`
}

type RefreshRes struct {
    AccessToken string `json:"access_token" dc:"新的访问令牌"`
}
```

3. Update LogoutRes (add cookie clearing note):
```go
type LogoutRes struct{}
```

- [ ] **Step 2: Commit**

```bash
git add api/v1/auth.go
git commit -m "feat: update auth API types for refresh token flow"
```

---

## Task 5: Backend — Update Auth Controller

**Files:**
- Modify: `internal/controller/auth.go`

- [ ] **Step 1: Update Login to return access_token and set refresh cookie**

Read the current file first. Make these changes:

1. In `Login`, rename the response field from `Token` to `AccessToken`:
```go
return &v1.LoginRes{
    AccessToken:  token,  // was Token: token,
    UserId:       user.Id,
    // ... rest same
}, nil
```

2. After generating the access token, also generate refresh token and set cookie:
```go
// After token, err := j.GenerateToken(...)
refreshToken, err := j.GenerateRefreshToken(user.Id)
if err != nil {
    return nil, gerror.New("生成refresh token失败")
}

// Set refresh token as httpOnly cookie
r := ghttp.RequestFromCtx(ctx)
r.Response.Cookie().Set(ghttp.CookieOption{
    Name:     "refresh_token",
    Value:    refreshToken,
    Path:     "/api/v1",
    HttpOnly: true,
    MaxAge:   604800, // 7 days
})
```

Note: You'll need to add `"github.com/gogf/gf/v2/net/ghttp"` to imports.

3. Add `Refresh` method:
```go
func (c *cAuth) Refresh(ctx context.Context, req *v1.RefreshReq) (*v1.RefreshRes, error) {
    r := ghttp.RequestFromCtx(ctx)
    refreshToken := r.Cookie("refresh_token")
    if refreshToken == "" {
        return nil, gerror.New("refresh token不存在")
    }

    jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
    secret := jwtConfig["secret"].String()
    if secret == "" {
        secret = "tool-go-jwt-secret-key-change-in-production"
    }
    expires := jwtConfig["expires"].Duration()
    if expires == 0 {
        expires = 15 * time.Minute
    }
    issuer := jwtConfig["issuer"].String()
    if issuer == "" {
        issuer = "tool-go"
    }

    j := jwt.New(secret, expires, issuer)
    claims, err := j.ParseRefreshToken(refreshToken)
    if err != nil {
        return nil, gerror.New("refresh token无效或已过期")
    }

    // Get user info for new tokens
    var user *entity.User
    err = dao.User.Ctx(ctx).
        Where(dao.User.Columns.Id, claims.UserId).
        WhereNull(dao.User.Columns.DeletedAt).
        Scan(&user)
    if err != nil || user == nil {
        return nil, gerror.New("用户不存在")
    }

    roles := getUserRoles(ctx, user.Id)
    permissions := getUserPermissions(ctx, user.Id)

    // Generate new access token
    newAccessToken, err := j.GenerateToken(user.Id, user.Username, roles, permissions)
    if err != nil {
        return nil, gerror.New("生成token失败")
    }

    // Generate new refresh token (rotation)
    newRefreshToken, err := j.GenerateRefreshToken(user.Id)
    if err != nil {
        return nil, gerror.New("生成refresh token失败")
    }

    // Set new refresh token cookie
    r.Response.Cookie().Set(ghttp.CookieOption{
        Name:     "refresh_token",
        Value:    newRefreshToken,
        Path:     "/api/v1",
        HttpOnly: true,
        MaxAge:   604800,
    })

    return &v1.RefreshRes{
        AccessToken: newAccessToken,
    }, nil
}
```

4. Update `Logout` to clear the cookie:
```go
func (c *cAuth) Logout(ctx context.Context, req *v1.LogoutReq) (*v1.LogoutRes, error) {
    r := ghttp.RequestFromCtx(ctx)
    r.Response.Cookie().Set(ghttp.CookieOption{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/api/v1",
        HttpOnly: true,
        MaxAge:   -1,
    })
    return &v1.LogoutRes{}, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/controller/auth.go
git commit -m "feat: implement refresh token flow in auth controller"
```

---

## Task 6: Backend — Wire Routes + Apply Middleware

**Files:**
- Modify: `internal/cmd/cmd.go`

- [ ] **Step 1: Add refresh route and security middleware**

Read the current file first.

1. Add `SecurityHeaders` to global middleware (line 107):
```go
group.Middleware(middleware.CORS, middleware.SecurityHeaders, ghttp.MiddlewareHandlerResponse)
```

2. Add rate limiter to login route (line 131):
```go
v1.POST("/login", middleware.RateLimit(5, time.Minute), controller.Auth, "Login")
```

Note: Add `"time"` to imports.

3. Add refresh route (after the login route, before the auth group):
```go
v1.POST("/refresh", controller.Auth, "Refresh")
```

- [ ] **Step 2: Commit**

```bash
git add internal/cmd/cmd.go
git commit -m "feat: wire refresh route, apply security headers and rate limiting"
```

---

## Task 7: Frontend — API and Store Updates

**Files:**
- Modify: `web/src/api/auth.ts`
- Modify: `web/src/store/modules/user.ts`

- [ ] **Step 1: Update auth.ts**

Read the current file first.

1. Update `LoginRes` to use `access_token`:
```typescript
export interface LoginRes {
  access_token: string  // was token
  user_id: number
  username: string
  nickname: string
  roles: string[]
  menus: MenuTree[]
  permissions: string[]
}
```

2. Add `RefreshRes`:
```typescript
export interface RefreshRes {
  access_token: string
}
```

3. Add `refresh` method to `authApi`:
```typescript
refresh: () => request.post<RefreshRes>('/refresh'),
```

- [ ] **Step 2: Update user.ts store**

Read the current file first.

Update `setToken` call site — the login response now has `access_token` instead of `token`. Find where `setToken` is called (likely in the login page or wherever login response is handled) and update the field name.

Also update the `setToken` initialization from localStorage — the key stays the same (`token`), but the source changes from `response.token` to `response.access_token`.

- [ ] **Step 3: Commit**

```bash
git add web/src/api/auth.ts web/src/store/modules/user.ts
git commit -m "feat: update frontend for refresh token flow"
```

---

## Task 8: Frontend — 401 Auto-Refresh Interceptor

**Files:**
- Modify: `web/src/utils/request.ts`

- [ ] **Step 1: Add 401 interceptor with refresh logic**

Read the current file first. In the response interceptor, add refresh logic:

```typescript
let isRefreshing = false
let pendingRequests: Array<(token: string) => void> = []

// In the response error interceptor:
if (error.response?.status === 401 && !originalRequest._retry) {
  if (isRefreshing) {
    return new Promise(resolve => {
      pendingRequests.push((token: string) => {
        originalRequest.headers.Authorization = `Bearer ${token}`
        resolve(request(originalRequest))
      })
    })
  }

  originalRequest._retry = true
  isRefreshing = true

  try {
    const { data } = await request.post('/refresh')
    const newToken = data.access_token
    localStorage.setItem('token', newToken)
    pendingRequests.forEach(cb => cb(newToken))
    pendingRequests = []
    originalRequest.headers.Authorization = `Bearer ${newToken}`
    return request(originalRequest)
  } catch {
    localStorage.removeItem('token')
    window.location.href = '/login'
    return Promise.reject(error)
  } finally {
    isRefreshing = false
  }
}
```

- [ ] **Step 2: Verify type check and build**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit 2>&1 && npx vite build 2>&1 | tail -5
```

- [ ] **Step 3: Commit**

```bash
git add web/src/utils/request.ts
git commit -m "feat: add 401 auto-refresh interceptor in Axios"
```

---

## Task 9: Frontend — Update Login Page Token Handling

**Files:**
- Modify: `web/src/views/login/index.vue` (or wherever login response is handled)

- [ ] **Step 1: Update login response handling**

Find where the login response is processed (likely `authApi.login()` call). Change `response.token` to `response.access_token`:

```typescript
// Before:
userStore.setToken(res.token)

// After:
userStore.setToken(res.access_token)
```

- [ ] **Step 2: Verify type check and build**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit 2>&1 && npx vite build 2>&1 | tail -5
```

- [ ] **Step 3: Commit**

```bash
git add web/src/views/login/index.vue
git commit -m "feat: update login page to use access_token field"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | Security headers middleware | 1 new |
| 2 | Rate limiter middleware | 1 new |
| 3 | Refresh token in JWT library | 1 modified |
| 4 | Auth API types update | 1 modified |
| 5 | Auth controller (Refresh/Logout) | 1 modified |
| 6 | Wire routes + apply middleware | 1 modified |
| 7 | Frontend API + store updates | 2 modified |
| 8 | 401 auto-refresh interceptor | 1 modified |
| 9 | Login page token handling | 1 modified |

**Total:** 2 new files, 9 modified files
