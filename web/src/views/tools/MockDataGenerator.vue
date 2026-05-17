<template>
  <div class="mock-data-generator">
    <div class="options">
      <!-- 数据类型选择区域 -->
      <div class="option-section">
        <div class="option-label">数据类型</div>
        <!-- el-checkbox-group：复选框组，支持多选，v-model 绑定选中的值数组 -->
        <el-checkbox-group v-model="selectedTypes" class="type-grid">
          <el-checkbox v-for="item in typeOptions" :key="item.value" :label="item.value" :value="item.value">
            {{ item.label }}
          </el-checkbox>
        </el-checkbox-group>
      </div>

      <!-- 生成数量控制区域 -->
      <div class="option-section">
        <div class="option-label">生成数量</div>
        <div class="count-control">
          <!-- el-slider：滑动条，直观调节数量（1～100） -->
          <el-slider v-model="count" :min="1" :max="100" style="flex: 1; margin-right: 16px" />
          <!-- el-input-number：数字输入框，可手动输入或步进调节 -->
          <el-input-number v-model="count" :min="1" :max="100" :step="1" size="small" style="width: 100px" />
        </div>
      </div>

      <!--
        :loading="loading"：按钮显示加载状态（旋转动画）
        :disabled="selectedTypes.length === 0"：未选择类型时禁用按钮
      -->
      <el-button type="primary" @click="generate" :loading="loading" :disabled="selectedTypes.length === 0">
        生成数据
      </el-button>
    </div>

    <!-- el-empty：未生成数据时的空状态提示 -->
    <el-empty v-if="!result && !loading" description="选择类型后点击「生成数据」" />

    <!-- v-if="result"：有结果时渲染表格区域 -->
    <div v-if="result" class="result-section">
      <div class="result-header">
        <span class="result-title">生成结果（共 {{ result.data.length }} 条）</span>
        <!-- 复制全部按钮：复制整个表格内容，用 tab 分隔列，换行分隔行 -->
        <el-button size="small" @click="copyAll">复制全部</el-button>
      </div>

      <div class="table-wrapper">
        <el-table :data="result.data" border stripe size="small" max-height="480">
          <!-- type="index" 自动生成行号，width="50" 固定列宽，fixed 固定在左侧 -->
          <el-table-column type="index" label="#" width="50" fixed />
          <!--
            v-for="col in result.columns"：动态渲染数据列
            :prop="col" 绑定数据字段名
            :label="columnLabel(col)" 显示中文列名
            :min-width="columnWidth(col)" 根据字段类型设置最小列宽，实现响应式
          -->
          <el-table-column
            v-for="col in result.columns"
            :key="col"
            :prop="col"
            :label="columnLabel(col)"
            :min-width="columnWidth(col)"
          >
            <!-- 默认插槽：通过 { row } 获取当前行数据 -->
            <template #default="{ row }">
              <span class="cell-text">{{ row[col] }}</span>
            </template>
          </el-table-column>
          <!-- 操作列：fixed="right" 固定在右侧 -->
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ $index }">
              <!-- 复制单行：将当前行所有字段以 "字段名: 值" 格式复制 -->
              <el-button size="small" type="primary" link @click="copyRow($index)">复制</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { toolsApi, MockDataRes } from '@/api/tools'

// 数据类型选项列表：label 是中文名，value 是字段 key
const typeOptions = [
  { label: '姓名', value: 'name' },
  { label: '手机号', value: 'phone' },
  { label: '邮箱', value: 'email' },
  { label: '身份证号', value: 'id_card' },
  { label: '护照号', value: 'passport' },
  { label: '地址', value: 'address' },
  { label: 'IP 地址', value: 'ip' },
  { label: '日期时间', value: 'datetime' },
  { label: '银行卡号', value: 'bank_card' },
]

// 英文 key → 中文列名的映射表，用于表头显示
const columnLabels: Record<string, string> = {
  name: '姓名',
  phone: '手机号',
  email: '邮箱',
  id_card: '身份证号',
  passport: '护照号',
  address: '地址',
  ip: 'IP 地址',
  datetime: '日期时间',
  bank_card: '银行卡号',
}

// 根据字段 key 获取中文列名，找不到则直接返回 key
const columnLabel = (key: string) => columnLabels[key] || key

// 根据字段 key 返回建议的最小列宽（像素），实现不同字段不同宽度
const columnWidth = (key: string) => {
  const widths: Record<string, number> = {
    id_card: 180,
    address: 240,
    passport: 140,
    bank_card: 160,
    email: 200,
    ip: 140,
    datetime: 160,
    phone: 130,
    name: 90,
  }
  return widths[key] || 120
}

// 用户选中的数据类型（复选框绑定的数组）
const selectedTypes = ref<string[]>([])
// 生成数量，默认 10
const count = ref(10)
// 加载状态，请求接口时设为 true，完成后设为 false
const loading = ref(false)
// 后端返回的生成结果（包含 columns 和 data），初始为 null
const result = ref<MockDataRes | null>(null)

// 调用后端 API 生成模拟数据
const generate = async () => {
  if (selectedTypes.value.length === 0) {
    ElMessage.warning('请至少选择一种数据类型')
    return
  }
  loading.value = true
  result.value = null
  try {
    // 向后端发送选中的类型和数量，获取模拟数据
    result.value = await toolsApi.mockData({ types: selectedTypes.value, count: count.value })
  } catch (err: any) {
    ElMessage.error(err.message || '生成失败')
  } finally {
    loading.value = false
  }
}

// 复制单行数据：格式为 "字段名1: 值1\n字段名2: 值2\n..."
const copyRow = (index: number) => {
  if (!result.value) return
  const row = result.value.data[index]
  const text = result.value.columns.map(k => `${columnLabel(k)}: ${row[k]}`).join('\n')
  navigator.clipboard.writeText(text).then(() => ElMessage.success('已复制'))
}

// 复制全部数据：第一行为列名（tab 分隔），后续每行为一行数据（tab 分隔）
const copyAll = () => {
  if (!result.value) return
  // 表头行：列名用 tab 连接
  const header = result.value.columns.map(k => columnLabel(k)).join('\t')
  // 数据行：每行的值用 tab 连接
  const rows = result.value.data.map(row =>
    result.value!.columns.map(k => row[k]).join('\t'),
  )
  // 表头 + 换行 + 数据行，整体用换行分隔
  const text = [header, ...rows].join('\n')
  navigator.clipboard.writeText(text).then(() => ElMessage.success('已复制全部'))
}
</script>

<style scoped>
.mock-data-generator {
  padding: 4px;

  .options {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 20px;
  }

  .option-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .option-label {
    font-size: 14px;
    font-weight: 500;
    color: #303133;
  }

  .type-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;

    :deep(.el-checkbox) {
      margin-right: 0;
      width: 120px;
    }
  }

  .count-control {
    display: flex;
    align-items: center;
  }

  .result-section {
    margin-top: 8px;

    .result-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 12px;

      .result-title {
        font-size: 14px;
        font-weight: 500;
        color: #303133;
      }
    }
  }

  .table-wrapper {
    overflow-x: auto;

    .cell-text {
      font-size: 13px;
      word-break: break-all;
    }
  }
}
</style>
