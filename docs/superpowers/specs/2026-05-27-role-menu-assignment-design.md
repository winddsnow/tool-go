# Role-Menu Assignment UI Design

> **Goal:** Add menu assignment capability to the role management page so admins can configure which menus each role can see.

**Architecture:** Add role-menu assignment API (get/set) and a tree-based checkbox dialog on the role management page.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend)

---

## 1. Backend: Role-Menu Assignment API

### Modified Files
- `api/v1/role.go` — Add GetMenus/AssignMenus structs
- `internal/service/role.go` — Extend IRole interface
- `internal/controller/role.go` — Add two methods
- `internal/logic/role.go` — Add two methods

### Endpoints
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/role/{id}/menus` | `Role.GetMenus` | Get role's menu IDs |
| PUT | `/role/{id}/menus` | `Role.AssignMenus` | Replace role's menus |

### API Structs
```go
type RoleGetMenusReq struct {
    g.Meta `path:"/role/{id}/menus" method:"get" tags:"Role" summary:"获取角色菜单"`
    Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

type RoleGetMenusRes struct {
    MenuIds []uint64 `json:"menu_ids" dc:"菜单ID列表"`
}

type RoleAssignMenusReq struct {
    g.Meta  `path:"/role/{id}/menus" method:"put" tags:"Role" summary:"分配角色菜单"`
    Id      uint64   `path:"id" v:"required#ID不能为空" dc:"角色ID"`
    MenuIds []uint64 `json:"menu_ids" v:"required#菜单列表不能为空" dc:"菜单ID列表"`
}

type RoleAssignMenusRes struct{}
```

### Logic
- `GetMenus`: Query `role_menu` table for `menu_id` values where `role_id = req.Id`
- `AssignMenus`: Transaction — DELETE all existing + INSERT new (same pattern as AssignPermissions)

---

## 2. Frontend: Role Page Menu Assignment

### Modified Files
- `web/src/api/role.ts` — Add getMenus/assignMenus methods
- `web/src/views/role/index.vue` — Add "分配菜单" button and tree dialog

### API Methods
```typescript
getMenus: (id: number) => request.get<{ menu_ids: number[] }>(`/role/${id}/menus`),
assignMenus: (id: number, data: { menu_ids: number[] }) =>
  request.put(`/role/${id}/menus`, data),
```

### UI Design

**Operations column:**
```
[编辑] [删除] [分配菜单] [分配权限]
```

**Menu assignment dialog:**
- `el-tree` with checkbox mode
- Shows all menus in tree hierarchy (parent → children)
- Supports parent-child cascade (checking parent checks all children)
- Loads current role's menu IDs on open
- Submit via `roleApi.assignMenus()`

### Menu Tree Building
- Fetch all menus via `menuApi.list({ page: 1, page_size: 1000 })`
- Build tree on frontend using `parent_id`
- Convert to `el-tree` data format with `id`, `label`, `children`

---

## Summary

| Component | Change |
|-----------|--------|
| Backend API | New role-menu assignment endpoints |
| Frontend API | New roleApi methods |
| Frontend UI | Role page "Assign Menu" tree dialog |
