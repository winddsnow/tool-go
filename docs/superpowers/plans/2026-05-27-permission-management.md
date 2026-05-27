# Permission Management Completion Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the RBAC permission management with API, UI, and auto-assignment for new users.

**Architecture:** Add permission list API, role-permission assignment API, auto-assign user role on creation, and "Assign Permissions" dialog on role management page.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend), PostgreSQL

---

## File Structure

### New Files
| File | Purpose |
|------|---------|
| `api/v1/permission.go` | Permission list API structs |
| `internal/controller/permission.go` | Permission controller |
| `internal/logic/permission.go` | Permission business logic |
| `internal/service/permission.go` | Permission service interface |
| `web/src/api/permission.ts` | Frontend permission API |

### Modified Files
| File | Change |
|------|--------|
| `api/v1/role.go` | Add GetPermissions/AssignPermissions structs |
| `internal/controller/role.go` | Add GetPermissions/AssignPermissions methods |
| `internal/logic/role.go` | Add GetPermissions/AssignPermissions logic |
| `internal/service/role.go` | Extend IRole interface |
| `internal/logic/user.go` | Auto-assign user role on creation |
| `internal/cmd/cmd.go` | Wire permission routes + role permission routes |
| `web/src/api/role.ts` | Add getPermissions/assignPermissions |
| `web/src/views/role/index.vue` | Add "Assign Permissions" button and dialog |

---

## Task 1: Backend — Permission Service + API

**Files:**
- Create: `internal/service/permission.go`
- Create: `api/v1/permission.go`
- Create: `internal/controller/permission.go`
- Create: `internal/logic/permission.go`

- [ ] **Step 1: Create service/permission.go**

```go
package service

import (
    "context"
    "sync"
    v1 "tool-go/api/v1"
)

type IPermission interface {
    List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error)
}

var (
    localPermission IPermission
    permissionMu    sync.RWMutex
)

func Permission() IPermission {
    permissionMu.RLock()
    defer permissionMu.RUnlock()
    return localPermission
}

func RegisterPermission(i IPermission) {
    permissionMu.Lock()
    defer permissionMu.Unlock()
    localPermission = i
}
```

- [ ] **Step 2: Create api/v1/permission.go**

```go
package v1

import "github.com/gogf/gf/v2/frame/g"

type PermissionListReq struct {
    g.Meta   `path:"/permission" method:"get" tags:"Permission" summary:"获取权限列表"`
    Page     int    `json:"page" d:"1" dc:"页码"`
    PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
    Code     string `json:"code" dc:"按权限码筛选"`
}

type PermissionListRes struct {
    List  []PermissionItem `json:"list" dc:"权限列表"`
    Total int              `json:"total" dc:"总数"`
    Page  int              `json:"page" dc:"当前页"`
}

type PermissionItem struct {
    Id     uint64 `json:"id" dc:"权限ID"`
    Code   string `json:"code" dc:"权限码"`
    Name   string `json:"name" dc:"权限名称"`
    MenuId uint64 `json:"menu_id" dc:"关联菜单ID"`
}
```

- [ ] **Step 3: Create controller/permission.go**

```go
package controller

import (
    "context"
    v1 "tool-go/api/v1"
    "tool-go/internal/service"
)

var Permission = cPermission{}

type cPermission struct{}

func (c *cPermission) List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error) {
    return service.Permission().List(ctx, req)
}
```

- [ ] **Step 4: Create logic/permission.go**

```go
package logic

import (
    "context"

    v1 "tool-go/api/v1"
    "tool-go/internal/dao"
    "tool-go/internal/model/entity"
    "tool-go/internal/service"
)

func init() {
    service.RegisterPermission(NewPermission())
}

func NewPermission() service.IPermission {
    return &sPermission{}
}

type sPermission struct{}

func (s *sPermission) List(ctx context.Context, req *v1.PermissionListReq) (*v1.PermissionListRes, error) {
    m := dao.Permission.Ctx(ctx)

    if req.Code != "" {
        m = m.WhereLike(dao.Permission.Columns.Code, "%"+req.Code+"%")
    }

    total, err := m.Count()
    if err != nil {
        return nil, err
    }

    var list []*entity.Permission
    err = m.Page(req.Page, req.PageSize).
        OrderDesc(dao.Permission.Columns.Id).
        Scan(&list)
    if err != nil {
        return nil, err
    }

    items := make([]v1.PermissionItem, 0, len(list))
    for _, p := range list {
        items = append(items, v1.PermissionItem{
            Id:     p.Id,
            Code:   p.Code,
            Name:   p.Name,
            MenuId: p.MenuId,
        })
    }

    return &v1.PermissionListRes{List: items, Total: total, Page: req.Page}, nil
}
```

- [ ] **Step 5: Commit**

```bash
git add internal/service/permission.go api/v1/permission.go internal/controller/permission.go internal/logic/permission.go
git commit -m "feat: add permission list API (service, controller, logic)"
```

---

## Task 2: Backend — Role Permission Assignment API

**Files:**
- Modify: `api/v1/role.go`
- Modify: `internal/service/role.go`
- Modify: `internal/controller/role.go`
- Modify: `internal/logic/role.go`

- [ ] **Step 1: Add API structs to api/v1/role.go**

Read the file first. Add at the end:

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

- [ ] **Step 2: Extend IRole interface in service/role.go**

Add two methods to the interface:
```go
type IRole interface {
    // ... existing methods ...
    GetPermissions(ctx context.Context, req *v1.RoleGetPermissionsReq) (*v1.RoleGetPermissionsRes, error)
    AssignPermissions(ctx context.Context, req *v1.RoleAssignPermissionsReq) error
}
```

- [ ] **Step 3: Add controller methods in controller/role.go**

```go
func (c *cRole) GetPermissions(ctx context.Context, req *v1.RoleGetPermissionsReq) (*v1.RoleGetPermissionsRes, error) {
    return service.Role().GetPermissions(ctx, req)
}

func (c *cRole) AssignPermissions(ctx context.Context, req *v1.RoleAssignPermissionsReq) error {
    return service.Role().AssignPermissions(ctx, req)
}
```

- [ ] **Step 4: Add logic methods in logic/role.go**

Read the file first. Add at the end:

```go
func (s *sRole) GetPermissions(ctx context.Context, req *v1.RoleGetPermissionsReq) (*v1.RoleGetPermissionsRes, error) {
    result, err := dao.RolePermission.Ctx(ctx).
        Where(dao.RolePermission.Columns.RoleId, req.Id).
        Fields(dao.RolePermission.Columns.PermissionId).
        Array()
    if err != nil {
        return nil, err
    }
    ids := make([]uint64, len(result))
    for i, v := range result {
        ids[i] = v.Uint64()
    }
    return &v1.RoleGetPermissionsRes{PermissionIds: ids}, nil
}

func (s *sRole) AssignPermissions(ctx context.Context, req *v1.RoleAssignPermissionsReq) error {
    return dao.RolePermission.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // Delete existing
        _, err := dao.RolePermission.Ctx(ctx).Where(dao.RolePermission.Columns.RoleId, req.Id).Delete()
        if err != nil {
            return err
        }
        // Insert new
        if len(req.PermissionIds) > 0 {
            data := make([]do.RolePermission, 0, len(req.PermissionIds))
            for _, pid := range req.PermissionIds {
                data = append(data, do.RolePermission{
                    RoleId:       req.Id,
                    PermissionId: pid,
                    CreatedAt:    gtime.Now(),
                })
            }
            _, err = dao.RolePermission.Ctx(ctx).Data(data).Insert()
            if err != nil {
                return err
            }
        }
        return nil
    })
}
```

Note: Add imports for `gdb`, `gtime`, and `do`.

- [ ] **Step 5: Commit**

```bash
git add api/v1/role.go internal/service/role.go internal/controller/role.go internal/logic/role.go
git commit -m "feat: add role permission assignment API (get/set)"
```

---

## Task 3: Backend — Auto-Assign user Role on User Creation

**Files:**
- Modify: `internal/logic/user.go`

- [ ] **Step 1: Update Create method**

Read the file first. In the `Create` method, after the user insert (after `result.LastInsertId()`), add:

```go
// 自动分配 user 角色 (role_id=3)
_, _ = dao.UserRole.Data(&do.UserRole{
    UserId: uint64(id),
    RoleId: 3,
}).Insert()
```

- [ ] **Step 2: Commit**

```bash
git add internal/logic/user.go
git commit -m "feat: auto-assign user role on user creation"
```

---

## Task 4: Backend — Wire Routes

**Files:**
- Modify: `internal/cmd/cmd.go`

- [ ] **Step 1: Add permission routes**

Read the file first. In the `auth` group, add `controller.Permission` to the `auth.Bind()` call:

```go
auth.Bind(
    controller.User,
    controller.Role,
    controller.Dashboard,
    controller.Menu,
    controller.Permission,  // ADD THIS
)
```

The role permission assignment routes (`/role/{id}/permissions`) will be automatically registered by `auth.Bind(controller.Role, ...)` since the API structs have `g.Meta` tags.

- [ ] **Step 2: Commit**

```bash
git add internal/cmd/cmd.go
git commit -m "feat: wire permission routes"
```

---

## Task 5: Frontend — Permission API

**Files:**
- Create: `web/src/api/permission.ts`
- Modify: `web/src/api/role.ts`

- [ ] **Step 1: Create web/src/api/permission.ts**

```typescript
import request from '@/utils/request'

export interface PermissionItem {
  id: number
  code: string
  name: string
  menu_id: number
}

export interface PermissionListRes {
  list: PermissionItem[]
  total: number
  page: number
}

export const permissionApi = {
  list: (params: { page: number; page_size: number; code?: string }) =>
    request.get<PermissionListRes>('/permission', { params }),
}
```

- [ ] **Step 2: Add methods to web/src/api/role.ts**

Read the file first. Add to the `roleApi` object:

```typescript
getPermissions: (id: number) => request.get<{ permission_ids: number[] }>(`/role/${id}/permissions`),
assignPermissions: (id: number, data: { permission_ids: number[] }) =>
  request.put(`/role/${id}/permissions`, data),
```

- [ ] **Step 3: Commit**

```bash
git add web/src/api/permission.ts web/src/api/role.ts
git commit -m "feat: add permission and role permission API"
```

---

## Task 6: Frontend — Role Page "Assign Permissions" Dialog

**Files:**
- Modify: `web/src/views/role/index.vue`

- [ ] **Step 1: Add permission dialog and button**

Read the current file first. Make these changes:

1. Add import at top:
```typescript
import { permissionApi, type PermissionItem } from '@/api/permission'
```

2. Add state variables:
```typescript
const permissionDialogVisible = ref(false)
const currentRoleId = ref(0)
const allPermissions = ref<PermissionItem[]>([])
const selectedPermissionIds = ref<number[]>([])
const permissionLoading = ref(false)
```

3. Add "分配权限" button in the operations column (after the existing edit/delete buttons):
```vue
<el-button type="warning" link @click="handleAssignPermissions(row)">分配权限</el-button>
```

4. Add methods:
```typescript
const handleAssignPermissions = async (row: RoleItem) => {
  currentRoleId.value = row.id
  permissionDialogVisible.value = true
  permissionLoading.value = true
  try {
    const [permRes, rolePermRes] = await Promise.all([
      permissionApi.list({ page: 1, page_size: 1000 }),
      roleApi.getPermissions(row.id),
    ])
    allPermissions.value = permRes.list || []
    selectedPermissionIds.value = rolePermRes.permission_ids || []
  } finally {
    permissionLoading.value = false
  }
}

const submitAssignPermissions = async () => {
  try {
    await roleApi.assignPermissions(currentRoleId.value, {
      permission_ids: selectedPermissionIds.value,
    })
    ElMessage.success('分配权限成功')
    permissionDialogVisible.value = false
  } catch {
    ElMessage.error('分配权限失败')
  }
}
```

5. Add dialog template (before the closing `</template>`):
```vue
<el-dialog v-model="permissionDialogVisible" title="分配权限" width="500px">
  <div v-loading="permissionLoading">
    <el-checkbox-group v-model="selectedPermissionIds">
      <el-checkbox
        v-for="perm in allPermissions"
        :key="perm.id"
        :label="perm.id"
        :value="perm.id"
      >
        {{ perm.name }} ({{ perm.code }})
      </el-checkbox>
    </el-checkbox-group>
  </div>
  <template #footer>
    <el-button @click="permissionDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="submitAssignPermissions">确定</el-button>
  </template>
</el-dialog>
```

- [ ] **Step 2: Verify type check and build**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit 2>&1 | head -20
```

- [ ] **Step 3: Commit**

```bash
git add web/src/views/role/index.vue
git commit -m "feat: add assign permissions dialog on role management page"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | Permission list API | 4 new |
| 2 | Role permission assignment API | 4 modified |
| 3 | Auto-assign user role | 1 modified |
| 4 | Wire routes | 1 modified |
| 5 | Frontend permission API | 1 new, 1 modified |
| 6 | Role page assign permissions dialog | 1 modified |

**Total:** 6 new files, 9 modified files
