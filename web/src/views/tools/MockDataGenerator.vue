<template>
  <div class="mock-data-generator">
    <div class="options">
      <div class="option-section">
        <div class="option-label">数据类型</div>
        <el-checkbox-group v-model="selectedTypes" class="type-grid">
          <el-checkbox v-for="item in typeOptions" :key="item.value" :label="item.value" :value="item.value">
            {{ item.label }}
          </el-checkbox>
        </el-checkbox-group>
      </div>

      <div class="option-section">
        <div class="option-label">生成数量</div>
        <div class="count-control">
          <el-slider v-model="count" :min="1" :max="100" style="flex: 1; margin-right: 16px" />
          <el-input-number v-model="count" :min="1" :max="100" :step="1" size="small" style="width: 100px" />
        </div>
      </div>

      <el-button type="primary" @click="generate" :loading="loading" :disabled="selectedTypes.length === 0">
        生成数据
      </el-button>
    </div>

    <el-empty v-if="!result && !loading" description="选择类型后点击「生成数据」" />

    <div v-if="result" class="result-section">
      <div class="result-header">
        <span class="result-title">生成结果（共 {{ result.data.length }} 条）</span>
        <el-button size="small" @click="copyAll">复制全部</el-button>
      </div>

      <div class="table-wrapper">
        <el-table :data="result.data" border stripe size="small" max-height="480">
          <el-table-column type="index" label="#" width="50" fixed />
          <el-table-column
            v-for="col in result.columns"
            :key="col"
            :prop="col"
            :label="columnLabel(col)"
            :min-width="columnWidth(col)"
          >
            <template #default="{ row }">
              <span class="cell-text">{{ row[col] }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ $index }">
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

const columnLabel = (key: string) => columnLabels[key] || key

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

const selectedTypes = ref<string[]>([])
const count = ref(10)
const loading = ref(false)
const result = ref<MockDataRes | null>(null)

const generate = async () => {
  if (selectedTypes.value.length === 0) {
    ElMessage.warning('请至少选择一种数据类型')
    return
  }
  loading.value = true
  result.value = null
  try {
    result.value = await toolsApi.mockData({ types: selectedTypes.value, count: count.value })
  } catch (err: any) {
    ElMessage.error(err.message || '生成失败')
  } finally {
    loading.value = false
  }
}

const copyRow = (index: number) => {
  if (!result.value) return
  const row = result.value.data[index]
  const text = result.value.columns.map(k => `${columnLabel(k)}: ${row[k]}`).join('\n')
  navigator.clipboard.writeText(text).then(() => ElMessage.success('已复制'))
}

const copyAll = () => {
  if (!result.value) return
  const header = result.value.columns.map(k => columnLabel(k)).join('\t')
  const rows = result.value.data.map(row =>
    result.value!.columns.map(k => row[k]).join('\t'),
  )
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
