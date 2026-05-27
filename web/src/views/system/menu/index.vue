<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" @click="handleAdd">新增菜单</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="searchForm" class="search-form" @submit.prevent="handleSearch">
        <el-form-item label="菜单名称">
          <el-input v-model="searchForm.name" placeholder="请输入菜单名称" clearable />
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
        <el-table
          :data="treeData"
          v-loading="loading"
          border
          stripe
          row-key="id"
          :tree-props="{ children: 'children' }"
          default-expand-all
        >
          <el-table-column prop="name" label="菜单名称" min-width="180" />
          <el-table-column prop="icon" label="图标" width="80" />
          <el-table-column prop="path" label="路由路径" min-width="140" />
          <el-table-column prop="component" label="组件路径" min-width="160" />
          <el-table-column prop="type" label="类型" width="90">
            <template #default="{ row }">
              <el-tag v-if="row.type === 1" type="warning">目录</el-tag>
              <el-tag v-else-if="row.type === 2" type="success">菜单</el-tag>
              <el-tag v-else type="info">按钮</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="70" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="visible" label="可见" width="80">
            <template #default="{ row }">
              <el-tag :type="row.visible === 1 ? 'success' : 'info'">
                {{ row.visible === 1 ? '显示' : '隐藏' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" class="responsive-dialog">
      <el-form :model="form" label-width="100px" class="menu-form">
        <el-form-item label="菜单名称">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="菜单类型">
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="目录" :value="1" />
            <el-option label="菜单" :value="2" />
            <el-option label="按钮" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentMenuOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            placeholder="请选择上级菜单（留空则为顶级）"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="路由路径">
          <el-input v-model="form.path" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item v-if="form.type === 2" label="组件路径">
          <el-input v-model="form.component" placeholder="请输入组件路径，如: system/menu/index" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="请输入图标名" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <el-form-item label="可见">
          <el-switch v-model="form.visible" :active-value="1" :inactive-value="0" active-text="显示" inactive-text="隐藏" />
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { menuApi, type MenuItem } from '@/api/menu'

const loading = ref(false)
const tableData = ref<MenuItem[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('')

const searchForm = reactive({
  name: '',
  status: undefined as number | undefined,
})

const form = reactive({
  id: 0,
  parent_id: 0,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  visible: 1,
  status: 1,
  type: 1,
})

const buildTree = (list: MenuItem[], parentId: number = 0): any[] => {
  return list
    .filter((item) => item.parent_id === parentId)
    .sort((a, b) => a.sort - b.sort)
    .map((item) => ({
      ...item,
      children: buildTree(list, item.id),
    }))
}

const treeData = computed(() => buildTree(tableData.value))

const parentMenuOptions = computed(() => {
  const all = [{ id: 0, name: '顶级菜单', children: [] as any[] }, ...treeData.value]
  return all
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await menuApi.list({
      page: 1,
      page_size: 1000,
      name: searchForm.name || undefined,
      status: searchForm.status,
    })
    tableData.value = res.list
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.status = undefined
  handleSearch()
}

const handleAdd = () => {
  dialogTitle.value = '新增菜单'
  Object.assign(form, { id: 0, parent_id: 0, name: '', path: '', component: '', icon: '', sort: 0, visible: 1, status: 1, type: 1 })
  dialogVisible.value = true
}

const handleEdit = (row: MenuItem) => {
  dialogTitle.value = '编辑菜单'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const handleDelete = async (row: MenuItem) => {
  await ElMessageBox.confirm('确定要删除该菜单吗？', '提示', { type: 'warning' })
  await menuApi.delete(row.id)
  ElMessage.success('删除成功')
  fetchData()
}

const handleSubmit = async () => {
  if (form.id) {
    await menuApi.update(form.id, {
      parent_id: form.parent_id,
      name: form.name,
      path: form.path,
      component: form.component,
      icon: form.icon,
      sort: form.sort,
      visible: form.visible,
      status: form.status,
      type: form.type,
    })
    ElMessage.success('更新成功')
  } else {
    await menuApi.create({
      parent_id: form.parent_id,
      name: form.name,
      path: form.path,
      component: form.component,
      icon: form.icon,
      sort: form.sort,
      visible: form.visible,
      status: form.status,
      type: form.type,
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  fetchData()
}

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

:deep(.responsive-dialog) {
  @media (max-width: 580px) {
    width: 95% !important;
  }

  .menu-form {
    @media (max-width: 580px) {
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
