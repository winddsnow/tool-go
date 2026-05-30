# Phase 2: Button-Level Permissions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add fine-grained button-level permissions so different roles see different UI elements and can access different API endpoints.

**Architecture:** Add `permission` + `role_permission` tables. JWT carries permission codes. Frontend `v-if` checks permission codes. Backend middleware validates permission codes from JWT.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend), PostgreSQL

---

## File Structure

### New Files
| File | Purpose |
|------|---------|
| `internal/model/entity/permission.go` | Entity struct for permission |
| `internal/model/entity/role_permission.go` | Entity struct for role_permission |
| `internal/model/do/permission.go` | DO struct for permission |
| `internal/model/do/role_permission.go` | DO struct for role_permission |
| `internal/dao/permission.go` | DAO for permission table |
| `internal/dao/role_permission.go` | DAO for role_permission table |
| `manifest/sql/20260527_add_permission_tables.sql` | Migration SQL |

### Modified Files
| File | Change |
|------|--------|
| `internal/library/jwt/jwt.go` | Add Permissions to Claims |
| `internal/middleware/auth.go` | Add GetPermissions/HasPermission/HasAnyPermission helpers |
| `internal/middleware/permission.go` | Add PermissionCode() middleware factory |
| `internal/controller/auth.go` | Include permissions in Login/GetUserInfo |
| `api/v1/auth.go` | Add Permissions field to LoginRes/GetUserInfoRes |
| `internal/cmd/cmd.go` | Replace Permission() with PermissionCode() |
| `web/src/store/modules/user.ts` | Add permissions state + hasPermission() |
| `web/src/api/auth.ts` | Add permissions to LoginRes/UserInfoRes |
| `web/src/layouts/default.vue` | Replace hasAnyRole with hasPermission |
| `web/src/views/user/index.vue` | Replace hasAnyRole with hasPermission |

---

## Task 1: Database Schema — Permission + Role_Permission Tables

**Files:**
- Create: `manifest/sql/20260527_add_permission_tables.sql`

- [ ] **Step 1: Create migration SQL**

```sql
-- Migration: Add permission + role_permission tables
-- Date: 2026-05-27

CREATE TABLE IF NOT EXISTS "permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "code" VARCHAR(64) NOT NULL,
    "name" VARCHAR(64) NOT NULL,
    "menu_id" BIGINT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_permission_code" ON "permission" ("code");

COMMENT ON TABLE "permission" IS '权限表';
COMMENT ON COLUMN "permission"."id" IS '权限ID';
COMMENT ON COLUMN "permission"."code" IS '权限码 (resource:action)';
COMMENT ON COLUMN "permission"."name" IS '权限名称';
COMMENT ON COLUMN "permission"."menu_id" IS '关联菜单ID (0=无关联)';

CREATE TABLE IF NOT EXISTS "role_permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "role_id" BIGINT NOT NULL,
    "permission_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_role_permission_role_perm" ON "role_permission" ("role_id", "permission_id");

COMMENT ON TABLE "role_permission" IS '角色权限关联表';

-- Seed permissions
INSERT INTO "permission" ("code", "name", "menu_id") VALUES
('user:create',       '创建用户',   3),
('user:delete',       '删除用户',   3),
('user:assign-roles', '分配角色',   3),
('role:create',       '创建角色',   4),
('role:delete',       '删除角色',   4),
('menu:create',       '创建菜单',   5),
('menu:delete',       '删除菜单',   5);

-- Seed role_permission
-- super_admin: all permissions
INSERT INTO "role_permission" ("role_id", "permission_id") VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7);

-- admin: user CRUD + role create/delete (no menu management)
INSERT INTO "role_permission" ("role_id", "permission_id") VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

-- user: no button-level permissions
```

- [ ] **Step 2: Also update init.sql**

Append the same DDL and seed data to `manifest/sql/init.sql`.

- [ ] **Step 3: Commit**

```bash
git add manifest/sql/20260527_add_permission_tables.sql manifest/sql/init.sql
git commit -m "feat: add permission + role_permission tables with seed data"
```

---

## Task 2: Backend — Entity, DO, DAO for Permission + Role_Permission

**Files:**
- Create: `internal/model/entity/permission.go`
- Create: `internal/model/entity/role_permission.go`
- Create: `internal/model/do/permission.go`
- Create: `internal/model/do/role_permission.go`
- Create: `internal/dao/permission.go`
- Create: `internal/dao/role_permission.go`

- [ ] **Step 1: Create entity.Permission**

```go
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Permission struct {
    Id        uint64      `orm:"id"`
    Code      string      `orm:"code"`
    Name      string      `orm:"name"`
    MenuId    uint64      `orm:"menu_id"`
    CreatedAt *gtime.Time `orm:"created_at"`
}
```

- [ ] **Step 2: Create entity.RolePermission**

```go
package entity

import "github.com/gogf/gf/v2/os/gtime"

type RolePermission struct {
    Id           uint64      `orm:"id"`
    RoleId       uint64      `orm:"role_id"`
    PermissionId uint64      `orm:"permission_id"`
    CreatedAt    *gtime.Time `orm:"created_at"`
}
```

- [ ] **Step 3: Create do.Permission**

```go
package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type Permission struct {
    g.Meta    `orm:"table:permission, do:true"`
    Id        any
    Code      any
    Name      any
    MenuId    any
    CreatedAt *gtime.Time
}
```

- [ ] **Step 4: Create do.RolePermission**

```go
package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type RolePermission struct {
    g.Meta       `orm:"table:role_permission, do:true"`
    Id           any
    RoleId       any
    PermissionId any
    CreatedAt    *gtime.Time
}
```

- [ ] **Step 5: Create dao.Permission**

```go
package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var Permission = permissionDao{}

type permissionDao struct {
    Table   string
    Group   string
    Columns permissionColumns
}

type permissionColumns struct {
    Id     string
    Code   string
    Name   string
    MenuId string
}

func init() {
    Permission = permissionDao{
        Table: "permission",
        Group: "default",
        Columns: permissionColumns{
            Id:     "id",
            Code:   "code",
            Name:   "name",
            MenuId: "menu_id",
        },
    }
}

func (d *permissionDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *permissionDao) Data(data *do.Permission) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
```

- [ ] **Step 6: Create dao.RolePermission**

```go
package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var RolePermission = rolePermissionDao{}

type rolePermissionDao struct {
    Table   string
    Group   string
    Columns rolePermissionColumns
}

type rolePermissionColumns struct {
    Id           string
    RoleId       string
    PermissionId string
}

func init() {
    RolePermission = rolePermissionDao{
        Table: "role_permission",
        Group: "default",
        Columns: rolePermissionColumns{
            Id:           "id",
            RoleId:       "role_id",
            PermissionId: "permission_id",
        },
    }
}

func (d *rolePermissionDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *rolePermissionDao) Data(data *do.RolePermission) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
```

- [ ] **Step 7: Commit**

```bash
git add internal/model/entity/permission.go internal/model/entity/role_permission.go \
        internal/model/do/permission.go internal/model/do/role_permission.go \
        internal/dao/permission.go internal/dao/role_permission.go
git commit -m "feat: add entity, DO, DAO for permission and role_permission tables"
```

---

## Task 3: Backend — Update JWT Claims + Auth Middleware Helpers

**Files:**
- Modify: `internal/library/jwt/jwt.go`
- Modify: `internal/middleware/auth.go`

- [ ] **Step 1: Add Permissions to JWT Claims**

In `internal/library/jwt/jwt.go`, add `Permissions` field to Claims struct:

```go
type Claims struct {
    jwt.RegisteredClaims
    UserId      uint64   `json:"user_id"`
    Username    string   `json:"username"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
}
```

Update `GenerateToken` signature to accept permissions:

```go
func (j *JWT) GenerateToken(userId uint64, username string, roles []string, permissions []string) (string, error) {
    claims := Claims{
        UserId:      userId,
        Username:    username,
        Roles:       roles,
        Permissions: permissions,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expires)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
        },
    }
    // ... rest of function
}
```

- [ ] **Step 2: Add permission helpers to auth middleware**

In `internal/middleware/auth.go`, add:

```go
const CtxPermissions = "permissions"

// In the Auth middleware, after injecting roles, also inject permissions:
ctx = context.WithValue(ctx, CtxPermissions, claims.Permissions)

// Add helper functions:
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

- [ ] **Step 3: Commit**

```bash
git add internal/library/jwt/jwt.go internal/middleware/auth.go
git commit -m "feat: add permissions to JWT claims and auth middleware helpers"
```

---

## Task 4: Backend — Add PermissionCode Middleware

**Files:**
- Modify: `internal/middleware/permission.go`

- [ ] **Step 1: Add PermissionCode middleware factory**

In `internal/middleware/permission.go`, add:

```go
// PermissionCode returns a middleware that checks if user has any of the required permission codes.
func PermissionCode(requiredPermissions ...string) func(r *ghttp.Request) {
    return func(r *ghttp.Request) {
        if !HasAnyPermission(r.GetCtx(), requiredPermissions...) {
            r.Response.WriteJsonExit(g.Map{
                "code":    403,
                "message": "没有权限访问",
            })
            return
        }
        r.Middleware.Next()
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/middleware/permission.go
git commit -m "feat: add PermissionCode middleware for button-level permission checks"
```

---

## Task 5: Backend — Update Login/UserInfo to Return Permissions

**Files:**
- Modify: `api/v1/auth.go`
- Modify: `internal/controller/auth.go`

- [ ] **Step 1: Add Permissions field to LoginRes and GetUserInfoRes**

In `api/v1/auth.go`:

```go
type LoginRes struct {
    Token       string     `json:"token" dc:"访问令牌"`
    UserId      uint64     `json:"user_id" dc:"用户ID"`
    Username    string     `json:"username" dc:"用户名"`
    Nickname    string     `json:"nickname" dc:"昵称"`
    Roles       []string   `json:"roles" dc:"角色列表"`
    Menus       []MenuTree `json:"menus" dc:"菜单树"`
    Permissions []string   `json:"permissions" dc:"权限码列表"`
}

type GetUserInfoRes struct {
    UserId      uint64     `json:"user_id" dc:"用户ID"`
    Username    string     `json:"username" dc:"用户名"`
    Nickname    string     `json:"nickname" dc:"昵称"`
    Roles       []string   `json:"roles" dc:"角色列表"`
    Menus       []MenuTree `json:"menus" dc:"菜单树"`
    Permissions []string   `json:"permissions" dc:"权限码列表"`
}
```

- [ ] **Step 2: Add getUserPermissions helper and update Login/GetUserInfo**

In `internal/controller/auth.go`, add helper function:

```go
func getUserPermissions(ctx context.Context, userId uint64) []string {
    result, err := dao.Permission.Ctx(ctx).
        LeftJoin("role_permission", "role_permission.permission_id=permission.id").
        LeftJoin("user_role", "user_role.role_id=role_permission.role_id").
        Where("user_role.user_id", userId).
        Fields(dao.Permission.Columns.Code).
        Array()
    if err != nil {
        g.Log().Error(ctx, "获取用户权限失败:", err)
        return []string{}
    }
    codes := make([]string, len(result))
    for i, v := range result {
        codes[i] = v.String()
    }
    return codes
}
```

Update `Login` method to include permissions:

```go
permissions := getUserPermissions(ctx, user.Id)
// ... generate token with permissions ...
token, err := j.GenerateToken(user.Id, user.Username, roles, permissions)
// ... return with permissions ...
return &v1.LoginRes{
    // ... existing fields ...
    Permissions: permissions,
}, nil
```

Update `GetUserInfo` similarly.

- [ ] **Step 3: Update GenerateToken call in Login**

The existing call is:
```go
token, err := j.GenerateToken(user.Id, user.Username, roles)
```

Change to:
```go
permissions := getUserPermissions(ctx, user.Id)
token, err := j.GenerateToken(user.Id, user.Username, roles, permissions)
```

- [ ] **Step 4: Commit**

```bash
git add api/v1/auth.go internal/controller/auth.go
git commit -m "feat: include permissions in login and getUserInfo responses"
```

---

## Task 6: Backend — Replace Route Protection with PermissionCode

**Files:**
- Modify: `internal/cmd/cmd.go`

- [ ] **Step 1: Replace Permission() with PermissionCode()**

In `internal/cmd/cmd.go`, change the three Permission middleware calls:

Line 195 (user group):
```go
// Before:
user.Middleware(middleware.Permission("super_admin", "admin"))
// After:
user.Middleware(middleware.PermissionCode("user:create", "user:delete", "user:assign-roles"))
```

Line 214 (role group):
```go
// Before:
role.Middleware(middleware.Permission("super_admin", "admin"))
// After:
role.Middleware(middleware.PermissionCode("role:create", "role:delete"))
```

Line 219 (menu group):
```go
// Before:
menu.Middleware(middleware.Permission("super_admin", "admin"))
// After:
menu.Middleware(middleware.PermissionCode("menu:create", "menu:delete"))
```

- [ ] **Step 2: Commit**

```bash
git add internal/cmd/cmd.go
git commit -m "feat: replace role-based route protection with permission-code checks"
```

---

## Task 7: Frontend — Store Update with Permissions

**Files:**
- Modify: `web/src/store/modules/user.ts`

- [ ] **Step 1: Add permissions state and helper methods**

Read the current store file, then add:

1. Import `ref` (already imported)
2. Add state: `const permissions = ref<string[]>(JSON.parse(localStorage.getItem('permissions') || '[]'))`
3. Add methods:
```typescript
function hasPermission(code: string): boolean {
  return permissions.value.includes(code)
}

function hasAnyPermission(codes: string[]): boolean {
  return codes.some(code => permissions.value.includes(code))
}
```
4. Update `setUserInfo` to accept `permissions?: string[]` and store it
5. Update `logout` to clear permissions
6. Export `permissions`, `hasPermission`, `hasAnyPermission`

- [ ] **Step 2: Commit**

```bash
git add web/src/store/modules/user.ts
git commit -m "feat: add permissions state and hasPermission helper to store"
```

---

## Task 8: Frontend — API Type Updates

**Files:**
- Modify: `web/src/api/auth.ts`

- [ ] **Step 1: Add permissions to LoginRes and UserInfoRes**

```typescript
export interface LoginRes {
  token: string
  user_id: number
  username: string
  nickname: string
  roles: string[]
  menus: MenuTree[]
  permissions: string[]   // ADD THIS
}

export interface UserInfoRes {
  user_id: number
  username: string
  nickname: string
  roles: string[]
  menus: MenuTree[]
  permissions: string[]   // ADD THIS
}
```

- [ ] **Step 2: Update router/index.ts to pass permissions**

In `web/src/router/index.ts`, after `userStore.setUserInfo(...)`, add permissions:

```typescript
userStore.setUserInfo({
  userId: userInfo.user_id,
  username: userInfo.username,
  nickname: userInfo.nickname,
  roles: userInfo.roles,
  menus: userInfo.menus,
  permissions: userInfo.permissions,  // ADD THIS
})
```

- [ ] **Step 3: Commit**

```bash
git add web/src/api/auth.ts web/src/router/index.ts
git commit -m "feat: add permissions to auth API types and router guard"
```

---

## Task 9: Frontend — Replace v-if Role Checks with Permission Checks

**Files:**
- Modify: `web/src/layouts/default.vue`
- Modify: `web/src/views/user/index.vue`

- [ ] **Step 1: Update default.vue header button**

Replace:
```vue
<el-button v-if="userStore.hasAnyRole(['super_admin', 'admin'])" type="primary" link @click="router.push('/user')">
```

With:
```vue
<el-button v-if="userStore.hasPermission('user:create')" type="primary" link @click="router.push('/user')">
```

- [ ] **Step 2: Update user/index.vue "新增用户" button**

Replace:
```vue
<el-button v-if="userStore.hasAnyRole(['super_admin', 'admin'])" type="primary" @click="handleAdd">新增用户</el-button>
```

With:
```vue
<el-button v-if="userStore.hasPermission('user:create')" type="primary" @click="handleAdd">新增用户</el-button>
```

- [ ] **Step 3: Verify type check and build**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit 2>&1 && npx vite build 2>&1 | tail -5
```

- [ ] **Step 4: Commit**

```bash
git add web/src/layouts/default.vue web/src/views/user/index.vue
git commit -m "feat: replace role-based v-if checks with permission-code checks"
```

---

## Summary

| Task | Description | Files Changed |
|------|-------------|---------------|
| 1 | Database schema | 2 SQL files |
| 2 | Entity, DO, DAO | 6 new Go files |
| 3 | JWT claims + auth middleware helpers | 2 modified Go files |
| 4 | PermissionCode middleware | 1 modified Go file |
| 5 | Login/UserInfo return permissions | 2 modified Go files |
| 6 | Replace route protection | 1 modified Go file |
| 7 | Frontend store permissions | 1 modified TS file |
| 8 | API type updates | 2 modified TS files |
| 9 | Replace v-if role checks | 2 modified Vue files |

**Total:** ~6 new files, ~13 modified files
