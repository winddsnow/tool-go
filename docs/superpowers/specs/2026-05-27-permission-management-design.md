# Permission Management Completion Design

> **Goal:** Complete the RBAC permission management system with API, UI, and auto-assignment for new users.

**Architecture:** Add permission list API, role-permission assignment API, auto-assign `user` role on user creation, and "Assign Permissions" button on role management page.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend), PostgreSQL

---

## 1. Backend: Permission List API

### New Files
- `api/v1/permission.go` — Request/response structs
- `internal/controller/permission.go` — Controller
- `internal/logic/permission.go` — Business logic
- `internal/service/permission.go` — Service interface

### Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/permission` | `Permission.List` | List all permissions (paginated) |

### API Structs
```go
// PermissionListReq
type PermissionListReq struct {
    g.Meta   `path:"/permission" method:"get" tags:"Permission" summary:"获取权限列表"`
    Page     int    `json:"page" d:"1" dc:"页码"`
    PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
    Code     string `json:"code" dc:"按权限码筛选"`
}

// PermissionListRes
type PermissionListRes struct {
    List  []PermissionItem `json:"list" dc:"权限列表"`
    Total int              `json:"total" dc:"总数"`
    Page  int              `json:"page" dc:"当前页"`
}

// PermissionItem
type PermissionItem struct {
    Id     uint64 `json:"id" dc:"权限ID"`
    Code   string `json:"code" dc:"权限码"`
    Name   string `json:"name" dc:"权限名称"`
    MenuId uint64 `json:"menu_id" dc:"关联菜单ID"`
}
```

---

## 2. Backend: Role-Permission Assignment API

### Modified Files
- `api/v1/role.go` — Add GetPermissions/AssignPermissions structs
- `internal/controller/role.go` — Add two methods
- `internal/logic/role.go` — Add two methods
- `internal/service/role.go` — Extend interface

### Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/role/{id}/permissions` | `Role.GetPermissions` | Get role's permission IDs |
| PUT | `/role/{id}/permissions` | `Role.AssignPermissions` | Replace role's permissions |

### API Structs
```go
type RoleGetPermissionsReq struct {
    g.Meta `path:"/role/{id}/permissions" method:"get" tags:"Role" summary:"获取角色权限"`
    Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

type RoleGetPermissionsRes struct {
    PermissionIds []uint64 `json:"permission_ids" dc:"权限ID列表"`
}

type RoleAssignPermissionsReq struct {
    g.Meta        `path:"/role/{id}/permissions" method:"put" tags:"Role" summary:"分配角色权限"`
    Id            uint64   `path:"id" v:"required#ID不能为空" dc:"角色ID"`
    PermissionIds []uint64 `json:"permission_ids" v:"required#权限列表不能为空" dc:"权限ID列表"`
}

type RoleAssignPermissionsRes struct{}
```

### Logic
- `GetPermissions`: Query `role_permission` table for `permission_id` values where `role_id = req.Id`
- `AssignRoles`: Transaction — DELETE all existing + INSERT new (same pattern as user/roles)

---

## 3. Backend: Auto-Assign user Role on User Creation

### Modified File
- `internal/logic/user.go` — Update `Create` method

After successful user insert, add:
```go
// 自动分配 user 角色 (role_id=3)
_, _ = dao.UserRole.Data(&do.UserRole{
    UserId: uint64(id),
    RoleId: 3,
}).Insert()
```

---

## 4. Backend: Wire Routes

### Modified File
- `internal/cmd/cmd.go`

Add permission routes and role permission routes:
```go
// Permission list (requires authentication)
auth.Bind(controller.Permission)

// Role permission assignment (requires super_admin or admin)
auth.Group("/role", func(role *ghttp.RouterGroup) {
    role.Middleware(middleware.PermissionCode("role:create", "role:delete"))
})
```

---

## 5. Frontend: Permission API

### New File
- `web/src/api/permission.ts`

```typescript
export const permissionApi = {
  list: (params: { page: number; page_size: number; code?: string }) =>
    request.get('/permission', { params }),
}
```

### Modified File
- `web/src/api/role.ts` — Add methods

```typescript
getPermissions: (id: number) => request.get(`/role/${id}/permissions`),
assignPermissions: (id: number, data: { permission_ids: number[] }) =>
  request.put(`/role/${id}/permissions`, data),
```

---

## 6. Frontend: Role Page "Assign Permissions" Button

### Modified File
- `web/src/views/role/index.vue`

Changes:
1. Add "分配权限" button in operations column (next to edit/delete)
2. Add permission assignment dialog with checkbox group
3. Load all permissions on dialog open
4. Load current role's permissions
5. Submit via `roleApi.assignPermissions()`

---

## Summary

| Component | Change |
|-----------|--------|
| Backend API | New permission list endpoint |
| Backend API | New role permission assignment endpoints |
| Backend Logic | Auto-assign user role on creation |
| Backend Routes | Wire new endpoints |
| Frontend API | New permissionApi + roleApi methods |
| Frontend UI | Role page "Assign Permissions" dialog |
