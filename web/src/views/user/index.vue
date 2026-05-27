<template>
  <div class="user-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button v-if="userStore.hasPermission('user:create')" type="primary" @click="handleAdd">新增用户</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form" @submit.prevent="handleSearch">
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <div class="table-wrapper">
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link @click="handleAssignRoles(row)">分配权限</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      </div>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @current-change="fetchData"
        @size-change="fetchData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" class="responsive-dialog">
      <el-form :model="form" label-width="100px" class="user-form">
        <el-form-item label="用户名" v-if="!form.id">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" v-if="!form.id">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialogVisible" title="分配权限" width="500px" class="responsive-dialog">
      <el-form label-width="100px">
        <el-form-item label="角色">
          <el-select v-model="selectedRoleIds" multiple placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="role in allRoles"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssignRoles">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
// ============================================================
// Vue 3 组合式 API —— 本文件是用户管理页面（CRUD）
// ref      ：响应式的基本类型或数组
// reactive ：响应式对象（深层响应）
// onMounted：生命周期钩子，组件挂载到 DOM 后自动执行
// ============================================================
import { ref, reactive, onMounted } from 'vue'

// ElMessage：全局消息提示（成功/失败）
// ElMessageBox：弹窗确认框（如"确定删除？"）
import { ElMessage, ElMessageBox } from 'element-plus'

// Pinia store：获取当前用户的角色信息（控制按钮显隐）
import { useUserStore } from '@/store/modules/user'

// userApi：用户相关 API（增删改查）
// UserItem：用户数据的 TypeScript 类型定义
import { userApi, UserItem } from '@/api/user'

// roleApi：角色相关 API
// RoleItem：角色数据的 TypeScript 类型定义
import { roleApi, RoleItem } from '@/api/role'

const userStore = useUserStore()

// loading：表格加载状态（v-loading 指令绑定，请求时显示加载动画）
const loading = ref(false)

// tableData：用户列表数据（响应式数组），供 el-table 渲染
const tableData = ref<UserItem[]>([])

// dialogVisible：新增/编辑弹窗的显示/隐藏（v-model 绑定）
// dialogTitle：弹窗标题（新增用户 / 编辑用户）
const dialogVisible = ref(false)
const dialogTitle = ref('')

// ---------- 角色分配弹窗相关 ----------
// roleDialogVisible：分配角色弹窗的显示/隐藏
// allRoles：所有可选的角色列表
// selectedRoleIds：当前用户已选中的角色 ID 数组
// currentUserId：正在操作的用户 ID
const roleDialogVisible = ref(false)
const allRoles = ref<RoleItem[]>([])
const selectedRoleIds = ref<number[]>([])
const currentUserId = ref(0)

// ----------------------------------------------------------
// searchForm：搜索表单数据模型（reactive 对象）
// v-model 双向绑定输入框的值
// status 是 number | undefined 可选值（全部/启用/禁用）
// ----------------------------------------------------------
const searchForm = reactive({
  username: '',
  status: undefined as number | undefined,
})

// ----------------------------------------------------------
// pagination：分页参数
// page     ：当前第几页
// pageSize ：每页显示条数
// total    ：总记录数（由后端返回）
// ----------------------------------------------------------
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// ----------------------------------------------------------
// form：新增/编辑弹窗中的表单数据模型
// id = 0 表示"新增"，id > 0 表示"编辑"
// Object.assign() 用来重置或填充表单数据
// ----------------------------------------------------------
const form = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  status: 1,
})

// ----------------------------------------------------------
// fetchData：获取用户列表（分页 + 搜索）
// async/await：等待后端 API 返回后继续执行
// finally：无论成功或失败都关闭 loading
// ----------------------------------------------------------
const fetchData = async () => {
  loading.value = true
  try {
    const res = await userApi.list({
      page: pagination.page,
      page_size: pagination.pageSize,
      username: searchForm.username || undefined,
      status: searchForm.status,
    })
    tableData.value = res.list   // 更新表格数据
    pagination.total = res.total // 更新总条数
  } finally {
    loading.value = false
  }
}

// fetchAllRoles：获取所有角色列表（用于分配角色弹窗中的下拉框）
const fetchAllRoles = async () => {
  try {
    const res = await roleApi.list({ page: 1, page_size: 1000 })
    allRoles.value = res.list
  } catch {
    // 失败时忽略（下拉框为空）
  }
}

// handleSearch：点击"查询"按钮
// 重置到第一页，然后重新请求数据
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// handleReset：点击"重置"按钮
// 清空搜索条件后重新搜索
const handleReset = () => {
  searchForm.username = ''
  searchForm.status = undefined
  handleSearch()
}

// ----------------------------------------------------------
// handleAdd：点击"新增用户"按钮
// Object.assign() 将表单重置为默认值（id=0 表示新增模式）
// 设置弹窗标题为"新增用户"，然后打开弹窗
// ----------------------------------------------------------
const handleAdd = () => {
  dialogTitle.value = '新增用户'
  Object.assign(form, { id: 0, username: '', password: '', nickname: '', email: '', phone: '', status: 1 })
  dialogVisible.value = true
}

// handleEdit：点击"编辑"按钮
// 将当前行的数据填充到 form 中（id > 0 表示编辑模式）
// 密码清空（编辑时不展示密码）
const handleEdit = (row: UserItem) => {
  dialogTitle.value = '编辑用户'
  Object.assign(form, { ...row, password: '' })
  dialogVisible.value = true
}

// ----------------------------------------------------------
// handleDelete：点击"删除"按钮
// ElMessageBox.confirm：弹出确认框，用户确认后才执行删除
// ----------------------------------------------------------
const handleDelete = async (row: UserItem) => {
  await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' })
  await userApi.delete(row.id)
  ElMessage.success('删除成功')
  fetchData() // 刷新列表
}

// ----------------------------------------------------------
// handleSubmit：弹窗中点击"确定"（新增或保存编辑）
// 根据 form.id 判断是新增（id=0）还是编辑（id>0）
// ----------------------------------------------------------
const handleSubmit = async () => {
  if (form.id) {
    // 编辑模式：调用更新接口（不传密码）
    await userApi.update(form.id, {
      username: form.username,
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      status: form.status,
    })
    ElMessage.success('更新成功')
  } else {
    // 新增模式：调用创建接口
    await userApi.create({
      username: form.username,
      password: form.password,
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      status: form.status,
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false // 关闭弹窗
  fetchData() // 刷新列表
}

// ----------------------------------------------------------
// handleAssignRoles：点击"分配权限"按钮
// 打开角色分配弹窗，并请求当前用户已拥有的角色
// ----------------------------------------------------------
const handleAssignRoles = async (row: UserItem) => {
  currentUserId.value = row.id
  roleDialogVisible.value = true
  try {
    const res = await userApi.getRoles(row.id)
    selectedRoleIds.value = res.role_ids
  } catch {
    selectedRoleIds.value = []
  }
}

// submitAssignRoles：角色分配弹窗中点击"确定"
// 将选中的角色 ID 提交到后端
const submitAssignRoles = async () => {
  try {
    await userApi.assignRoles(currentUserId.value, selectedRoleIds.value)
    ElMessage.success('分配成功')
    roleDialogVisible.value = false
  } catch {
    // 失败时忽略
  }
}

// ----------------------------------------------------------
// onMounted：组件挂载到页面时自动执行
// 等同于 Vue 2 的 mounted() 生命周期
// 页面一加载就请求用户列表和角色列表
// ----------------------------------------------------------
onMounted(() => {
  fetchData()
  fetchAllRoles()
})
</script>

<style scoped lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.search-form {
  margin-bottom: 20px;

  @media (max-width: 480px) {
    :deep(.el-form-item) {
      display: block;
      margin-right: 0;
      margin-bottom: 8px;
    }
    :deep(.el-form-item__content) {
      display: block;
    }
  }
}

.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

:deep(.responsive-dialog) {
  @media (max-width: 520px) {
    width: 95% !important;
  }

  .user-form {
    @media (max-width: 520px) {
      :deep(.el-form-item) {
        display: block;
        margin-bottom: 12px;
      }
      :deep(.el-form-item__label) {
        float: none;
        display: block;
        text-align: left;
        padding: 0 0 4px;
      }
      :deep(.el-form-item__content) {
        margin-left: 0 !important;
      }
    }
  }
}
</style>
