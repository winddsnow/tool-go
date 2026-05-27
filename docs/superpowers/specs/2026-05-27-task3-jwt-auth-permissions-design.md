# Task 3: JWT Claims + Auth Middleware Permissions Design

## Objective
Add permissions support to JWT tokens and auth middleware helpers for fine-grained button-level permissions.

## Scope
1. JWT Claims: Add `Permissions []string` field
2. JWT GenerateToken: Accept permissions parameter
3. Auth middleware: Add `CtxPermissions` constant and inject permissions into context
4. Auth middleware helpers: Add `GetPermissions`, `HasPermission`, `HasAnyPermission`

**Out of scope:** PermissionCode middleware (Task 4), frontend changes, login response updates.

## Changes

### 1. `internal/library/jwt/jwt.go`

**Claims struct modification:**
```go
type Claims struct {
    UserId      uint64   `json:"user_id"`
    Username    string   `json:"username"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}
```

**GenerateToken signature update:**
```go
func (j *JWT) GenerateToken(userId uint64, username string, roles []string, permissions []string) (string, error) {
```

**Claims initialization:**
```go
claims := Claims{
    UserId:      userId,
    Username:    username,
    Roles:       roles,
    Permissions: permissions,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(j.expires)),
        IssuedAt:  jwt.NewNumericDate(now),
        NotBefore: jwt.NewNumericDate(now),
        Issuer:    j.issuer,
    },
}
```

### 2. `internal/middleware/auth.go`

**Add context key constant:**
```go
const (
    CtxUserId      = "userId"
    CtxUsername    = "username"
    CtxRoles       = "roles"
    CtxPermissions = "permissions"
)
```

**Inject permissions in Auth middleware (after line 105):**
```go
ctx = context.WithValue(ctx, CtxPermissions, claims.Permissions)
```

**Add helper functions (after HasAnyRole):**
```go
func GetPermissions(ctx context.Context) []string {
    permissions, _ := ctx.Value(CtxPermissions).([]string)
    return permissions
}

func HasPermission(ctx context.Context, permission string) bool {
    for _, p := range GetPermissions(ctx) {
        if p == permission {
            return true
        }
    }
    return false
}

func HasAnyPermission(ctx context.Context, permissions ...string) bool {
    userPerms := GetPermissions(ctx)
    for _, required := range permissions {
        for _, p := range userPerms {
            if p == required {
                return true
            }
        }
    }
    return false
}
```

## Trade-offs
- **Simplicity:** Direct field addition to Claims struct, minimal changes.
- **Backward compatibility:** Existing tokens without permissions will decode with empty slice (Go zero value).
- **Performance:** No additional overhead; permissions carried in JWT payload.

## Verification
- After implementation, run `go build ./...` to ensure compilation.
- Check that existing callers of `GenerateToken` compile (they'll need to be updated separately).

## Design Approval
Please approve this design before I proceed with implementation.