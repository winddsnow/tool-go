<template>
  <div class="role-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">新增角色</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form" @submit.prevent="handleSearch">
        <el-form-item label="角色名称">
          <el-input v-model="searchForm.name" placeholder="请输入角色名称" clearable />
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
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="code" label="角色编码" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="desc" label="描述" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="角色名称">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码">
          <el-input v-model="form.code" placeholder="请输入角色编码" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.desc" type="textarea" placeholder="请输入描述" />
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
  </div>
</template>

<script setup lang="ts">
// ============================================================
// Vue 3 组合式 API —— 本文件是角色管理页面
// ref      ：响应式的基本类型或数组（用 .value 读写）
// reactive ：响应式对象（直接修改属性即可触发更新）
// onMounted：生命周期钩子，组件挂载后自动执行
// ============================================================
import { ref, reactive, onMounted } from 'vue'

// ElMessage：全局消息提示
// ElMessageBox：确认弹窗
import { ElMessage, ElMessageBox } from 'element-plus'

// roleApi：角色相关 API（增删改查）
// RoleItem：角色数据项的 TypeScript 类型
import { roleApi, RoleItem } from '@/api/role'

// loading：表格加载状态（绑定 v-loading，请求时显示遮罩动画）
const loading = ref(false)

// tableData：角色列表数据，供 el-table 渲染
const tableData = ref<RoleItem[]>([])

// dialogVisible：新增/编辑弹窗是否显示
// dialogTitle：弹窗标题文字
const dialogVisible = ref(false)
const dialogTitle = ref('')

// ---------- 搜索表单 ----------
// searchForm：搜索条件（角色名称、状态）
// status 为可选值（undefined 表示"全部"）
const searchForm = reactive({
  name: '',
  status: undefined as number | undefined,
})

// ---------- 分页参数 ----------
// page    ：当前页码
// pageSize：每页条数
// total   ：后端返回的总记录数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// ---------- 弹窗表单数据 ----------
// form：角色编辑表单，id=0 为新增，id>0 为编辑
// 与用户管理类似，但多了一个 "角色编码 (code)" 字段
const form = reactive({
  id: 0,
  name: '',
  code: '',
  sort: 0,
  desc: '',
  status: 1,
})

// fetchData：请求角色列表（分页 + 按名称/状态筛选）
const fetchData = async () => {
  loading.value = true
  try {
    const res = await roleApi.list({
      page: pagination.page,
      page_size: pagination.pageSize,
      name: searchForm.name || undefined,
      status: searchForm.status,
    })
    tableData.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

// handleSearch：点击"查询"（重置到第一页再搜索）
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// handleReset：点击"重置"（清空搜索条件）
const handleReset = () => {
  searchForm.name = ''
  searchForm.status = undefined
  handleSearch()
}

// handleAdd：新增角色（重置表单，打开弹窗）
const handleAdd = () => {
  dialogTitle.value = '新增角色'
  Object.assign(form, { id: 0, name: '', code: '', sort: 0, desc: '', status: 1 })
  dialogVisible.value = true
}

// handleEdit：编辑角色（将行数据填充到表单）
const handleEdit = (row: RoleItem) => {
  dialogTitle.value = '编辑角色'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

// handleDelete：删除角色（先弹确认框）
const handleDelete = async (row: RoleItem) => {
  await ElMessageBox.confirm('确定要删除该角色吗？', '提示', { type: 'warning' })
  await roleApi.delete(row.id)
  ElMessage.success('删除成功')
  fetchData()
}

// ----------------------------------------------------------
// handleSubmit：提交角色表单（新增/编辑）
// 与用户管理不同的是：
//   编辑时角色编码 (code) 不可修改
//   模板中的 :disabled="!!form.id" 在编辑模式下禁用该输入框
// ----------------------------------------------------------
const handleSubmit = async () => {
  if (form.id) {
    // 编辑模式
    await roleApi.update(form.id, {
      name: form.name,
      code: form.code,
      sort: form.sort,
      desc: form.desc,
      status: form.status,
    })
    ElMessage.success('更新成功')
  } else {
    // 新增模式
    await roleApi.create({
      name: form.name,
      code: form.code,
      sort: form.sort,
      desc: form.desc,
      status: form.status,
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  fetchData()
}

// onMounted：页面加载时自动拉取角色列表
onMounted(fetchData)
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
</style>
