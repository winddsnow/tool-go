# Role-Menu Assignment UI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add menu assignment capability to the role management page with tree-based checkbox dialog.

**Architecture:** Add role-menu assignment API (get/set) following the same pattern as role-permission. Add tree checkbox dialog on role management page.

**Tech Stack:** GoFrame v2 (backend), Vue 3 + TypeScript + Element Plus (frontend)

---

## File Structure

### Modified Files
| File | Change |
|------|--------|
| `api/v1/role.go` | Add GetMenus/AssignMenus structs |
| `internal/service/role.go` | Extend IRole interface |
| `internal/controller/role.go` | Add GetMenus/AssignMenus methods |
| `internal/logic/role.go` | Add GetMenus/AssignMenus logic |
| `web/src/api/role.ts` | Add getMenus/assignMenus methods |
| `web/src/views/role/index.vue` | Add "分配菜单" button and tree dialog |

---

## Task 1: Backend — Role-Menu Assignment API

**Files:**
- Modify: `api/v1/role.go`
- Modify: `internal/service/role.go`
- Modify: `internal/controller/role.go`
- Modify: `internal/logic/role.go`

- [ ] **Step 1: Add API structs to api/v1/role.go**

Read the file first. Add at the end (after the permission structs):

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

- [ ] **Step 2: Extend IRole interface in service/role.go**

Add two methods:
```go
GetMenus(ctx context.Context, req *v1.RoleGetMenusReq) (*v1.RoleGetMenusRes, error)
AssignMenus(ctx context.Context, req *v1.RoleAssignMenusReq) error
```

- [ ] **Step 3: Add controller methods in controller/role.go**

```go
func (c *cRole) GetMenus(ctx context.Context, req *v1.RoleGetMenusReq) (*v1.RoleGetMenusRes, error) {
    return service.Role().GetMenus(ctx, req)
}

func (c *cRole) AssignMenus(ctx context.Context, req *v1.RoleAssignMenusReq) error {
    return service.Role().AssignMenus(ctx, req)
}
```

- [ ] **Step 4: Add logic methods in logic/role.go**

Read the file first. Add at the end:

```go
func (s *sRole) GetMenus(ctx context.Context, req *v1.RoleGetMenusReq) (*v1.RoleGetMenusRes, error) {
    result, err := dao.RoleMenu.Ctx(ctx).
        Where(dao.RoleMenu.Columns.RoleId, req.Id).
        Fields(dao.RoleMenu.Columns.MenuId).
        Array()
    if err != nil {
        return nil, err
    }
    ids := make([]uint64, len(result))
    for i, v := range result {
        ids[i] = v.Uint64()
    }
    return &v1.RoleGetMenusRes{MenuIds: ids}, nil
}

func (s *sRole) AssignMenus(ctx context.Context, req *v1.RoleAssignMenusReq) error {
    return dao.RoleMenu.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        _, err := dao.RoleMenu.Ctx(ctx).Where(dao.RoleMenu.Columns.RoleId, req.Id).Delete()
        if err != nil {
            return err
        }
        if len(req.MenuIds) > 0 {
            data := make([]do.RoleMenu, 0, len(req.MenuIds))
            for _, mid := range req.MenuIds {
                data = append(data, do.RoleMenu{
                    RoleId:    req.Id,
                    MenuId:    mid,
                    CreatedAt: gtime.Now(),
                })
            }
            _, err = dao.RoleMenu.Ctx(ctx).Data(data).Insert()
            if err != nil {
                return err
            }
        }
        return nil
    })
}
```

Note: `gdb` and `gtime` imports should already be present from the permission methods added earlier. Make sure `"tool-go/internal/model/do"` is imported.

- [ ] **Step 5: Commit**

```bash
git add api/v1/role.go internal/service/role.go internal/controller/role.go internal/logic/role.go
git commit -m "feat: add role menu assignment API (get/set)"
```

---

## Task 2: Frontend — Role API Methods

**Files:**
- Modify: `web/src/api/role.ts`

- [ ] **Step 1: Add getMenus and assignMenus methods**

Read the file first. Add to the `roleApi` object:

```typescript
getMenus: (id: number) => request.get<{ menu_ids: number[] }>(`/role/${id}/menus`),
assignMenus: (id: number, data: { menu_ids: number[] }) =>
  request.put(`/role/${id}/menus`, data),
```

- [ ] **Step 2: Commit**

```bash
git add web/src/api/role.ts
git commit -m "feat: add role menu API methods"
```

---

## Task 3: Frontend — Role Page Menu Assignment Dialog

**Files:**
- Modify: `web/src/views/role/index.vue`

- [ ] **Step 1: Read the current file and understand the structure**

Read the file to understand:
- Where imports are
- Where state variables are defined
- Where the operations column buttons are
- Where methods are defined
- Where the template ends

- [ ] **Step 2: Add import for menuApi**

Add at the top of `<script setup>`:
```typescript
import { menuApi, type MenuTree } from '@/api/menu'
```

- [ ] **Step 3: Add state variables**

```typescript
const menuDialogVisible = ref(false)
const currentRoleIdForMenu = ref(0)
const menuTreeData = ref<MenuTree[]>([])
const selectedMenuIds = ref<number[]>([])
const menuLoading = ref(false)
```

- [ ] **Step 4: Add menu tree building helper**

```typescript
interface TreeMenuItem {
  id: number
  label: string
  children?: TreeMenuItem[]
}

function buildMenuTree(menus: MenuTree[]): TreeMenuItem[] {
  const map = new Map<number, TreeMenuItem>()
  const roots: TreeMenuItem[] = []

  menus.forEach(m => {
    map.set(m.id, { id: m.id, label: m.name, children: [] })
  })

  menus.forEach(m => {
    const node = map.get(m.id)!
    if (m.parent_id === 0) {
      roots.push(node)
    } else {
      const parent = map.get(m.parent_id)
      if (parent) {
        parent.children!.push(node)
      } else {
        roots.push(node)
      }
    }
  })

  return roots
}
```

- [ ] **Step 5: Add "分配菜单" button in operations column**

Find the operations column and add the button (before or after "分配权限"):
```vue
<el-button type="success" link @click="handleAssignMenus(row)">分配菜单</el-button>
```

- [ ] **Step 6: Add methods**

```typescript
const handleAssignMenus = async (row: RoleItem) => {
  currentRoleIdForMenu.value = row.id
  menuDialogVisible.value = true
  menuLoading.value = true
  try {
    const [menuRes, roleMenuRes] = await Promise.all([
      menuApi.getUserMenus(),
      roleApi.getMenus(row.id),
    ])
    menuTreeData.value = buildMenuTree(menuRes.menus || [])
    selectedMenuIds.value = roleMenuRes.menu_ids || []
  } finally {
    menuLoading.value = false
  }
}

const submitAssignMenus = async () => {
  try {
    await roleApi.assignMenus(currentRoleIdForMenu.value, {
      menu_ids: selectedMenuIds.value,
    })
    ElMessage.success('分配菜单成功')
    menuDialogVisible.value = false
  } catch {
    ElMessage.error('分配菜单失败')
  }
}
```

- [ ] **Step 7: Add dialog template (before closing `</template>`)**

```vue
<el-dialog v-model="menuDialogVisible" title="分配菜单" width="500px">
  <div v-loading="menuLoading">
    <el-tree
      :data="menuTreeData"
      show-checkbox
      node-key="id"
      :default-checked-keys="selectedMenuIds"
      :props="{ label: 'label', children: 'children' }"
      ref="menuTreeRef"
    />
  </div>
  <template #footer>
    <el-button @click="menuDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="handleMenuSubmit">确定</el-button>
  </template>
</el-dialog>
```

Note: The `el-tree` with `show-checkbox` and `ref` needs a submit handler that reads checked keys:

```typescript
import { ref } from 'vue'

const menuTreeRef = ref()

const handleMenuSubmit = async () => {
  if (!menuTreeRef.value) return
  const checkedKeys = menuTreeRef.value.getCheckedKeys()
  const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
  const allKeys = [...checkedKeys, ...halfCheckedKeys]
  
  try {
    await roleApi.assignMenus(currentRoleIdForMenu.value, {
      menu_ids: allKeys,
    })
    ElMessage.success('分配菜单成功')
    menuDialogVisible.value = false
  } catch {
    ElMessage.error('分配菜单失败')
  }
}
```

- [ ] **Step 8: Verify type check and build**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit 2>&1 | head -20
```

- [ ] **Step 9: Commit**

```bash
git add web/src/views/role/index.vue
git commit -m "feat: add assign menus dialog on role management page"
```

---

## Summary

| Task | Description | Files |
|------|-------------|-------|
| 1 | Backend role-menu assignment API | 4 modified |
| 2 | Frontend role API methods | 1 modified |
| 3 | Role page menu assignment dialog | 1 modified |

**Total:** 0 new files, 6 modified files
