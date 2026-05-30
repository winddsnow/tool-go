# Phase 2: Button-Level Permissions Design

> **Goal:** Add fine-grained button-level permissions so different roles see different UI elements and can access different API endpoints.

**Architecture:** Add `permission` + `role_permission` tables. JWT carries permission codes. Frontend `v-if` checks permission codes. Backend middleware validates permission codes from JWT.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend), PostgreSQL

---

## Data Model

### New Tables

```sql
CREATE TABLE IF NOT EXISTS "permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "code" VARCHAR(64) NOT NULL UNIQUE,
    "name" VARCHAR(64) NOT NULL,
    "menu_id" BIGINT DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "role_permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "role_id" BIGINT NOT NULL,
    "permission_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "idx_role_permission_role_perm" ON "role_permission" ("role_id", "permission_id");
```

### Permission Codes

| Code | Name | Description |
|------|------|-------------|
| `user:create` | 创建用户 | 新增用户按钮 |
| `user:delete` | 删除用户 | 删除用户按钮 |
| `user:assign-roles` | 分配角色 | 用户角色分配 |
| `role:create` | 创建角色 | 新增角色按钮 |
| `role:delete` | 删除角色 | 删除角色按钮 |
| `menu:create` | 创建菜单 | 新增菜单按钮 |
| `menu:delete` | 删除菜单 | 删除菜单按钮 |

### Seed Data: Role-Permission Mapping

| Role | Permissions |
|------|-------------|
| super_admin | ALL (7 permissions) |
| admin | user:create, user:delete, user:assign-roles, role:create, role:delete |
| user | (none) |

---

## Backend Changes

### 1. JWT Claims

Add `Permissions []string` to JWT claims:

```go
type Claims struct {
    jwt.RegisteredClaims
    UserId      uint64   `json:"user_id"`
    Username    string   `json:"username"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
}
```

### 2. Auth Middleware

Add `GetPermissions(ctx)` and `HasPermission(ctx, code)` helpers:

```go
func GetPermissions(ctx context.Context) []string { ... }
func HasPermission(ctx context.Context, code string) bool { ... }
```

### 3. Permission Middleware

New middleware factory:

```go
func PermissionCode(requiredPermissions ...string) func(r *ghttp.Request) {
    return func(r *ghttp.Request) {
        if !HasAnyPermission(r.GetCtx(), requiredPermissions...) {
            r.Response.WriteJsonExit(g.Map{"code": 403, "message": "没有权限"})
            return
        }
        r.Middleware.Next()
    }
}
```

### 4. Login Response

Include permissions in LoginRes and GetUserInfoRes:

```go
type LoginRes struct {
    ...
    Permissions []string `json:"permissions"`
}
```

### 5. Route Protection

Replace role-based middleware with permission-based:

```go
// Before:
user.Middleware(middleware.Permission("super_admin", "admin"))

// After:
user.Middleware(middleware.PermissionCode("user:create", "user:delete", "user:assign-roles"))
```

---

## Frontend Changes

### 1. Store

Add to `user.ts`:

```typescript
const permissions = ref<string[]>([])

function hasPermission(code: string): boolean {
  return permissions.value.includes(code)
}

function hasAnyPermission(codes: string[]): boolean {
  return codes.some(code => permissions.value.includes(code))
}
```

### 2. API Types

Update `LoginRes` and `UserInfoRes` to include `permissions: string[]`.

### 3. Template Checks

Replace:
```vue
<el-button v-if="userStore.hasAnyRole(['super_admin', 'admin'])" ...>
```

With:
```vue
<el-button v-if="userStore.hasPermission('user:create')" ...>
```

### 4. Router Guard

No changes needed — dynamic routes already rely on authentication, not role checks.

---

## Migration SQL

File: `manifest/sql/20260527_add_permission_tables.sql`

Contains CREATE TABLE, indexes, seed data for permissions and role_permission.

---

## Summary

| Component | Change |
|-----------|--------|
| DB | 2 new tables: permission, role_permission |
| JWT | Add permissions field to claims |
| Auth middleware | Add GetPermissions/HasPermission helpers |
| Permission middleware | Add PermissionCode() factory |
| Login/UserInfo | Return permissions list |
| Route protection | Replace Permission() with PermissionCode() |
| Frontend store | Add permissions state + hasPermission() |
| Frontend templates | Replace hasAnyRole() with hasPermission() |
| Seed data | 7 permissions, role mappings for 3 roles |
