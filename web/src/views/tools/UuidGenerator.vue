<template>
  <div class="uuid-generator">
    <div class="options">
      <div class="option-item">
        <span class="label">UUID 版本</span>
        <el-select v-model="version" size="large" style="width: 160px">
          <el-option label="UUID v4 (随机)" value="v4" />
          <el-option label="UUID v1 (时间戳)" value="v1" />
        </el-select>
      </div>
      <div class="option-item">
        <span class="label">生成数量</span>
        <el-input-number v-model="count" :min="1" :max="20" size="large" />
      </div>
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="generate">生成</el-button>
    </div>
    <div class="result" v-if="uuids.length > 0">
      <div class="uuid-item" v-for="(uuid, i) in uuids" :key="i">
        <span class="uuid-text">{{ uuid }}</span>
        <el-button size="small" link type="primary" @click="copy(uuid)">复制</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const version = ref('v4')
const count = ref(5)
const uuids = ref<string[]>([])

function generate() {
  uuids.value = []
  for (let i = 0; i < count.value; i++) {
    uuids.value.push(version.value === 'v4' ? generateV4() : generateV1())
  }
}

function generateV4(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  const hex = '0123456789abcdef'
  const sections = [8, 4, 4, 4, 12]
  return sections.map(len => {
    let s = ''
    for (let i = 0; i < len; i++) s += hex[Math.floor(Math.random() * 16)]
    return s
  }).join('-')
}

function generateV1(): string {
  const now = Date.now()
  const timeHex = now.toString(16).padStart(15, '0')
  const node = Math.floor(Math.random() * 0xffffffffffff).toString(16).padStart(12, '0')
  const clockSeq = Math.floor(Math.random() * 0x3fff).toString(16).padStart(4, '0')
  return `${timeHex.slice(0, 8)}-${timeHex.slice(8, 12)}-1${timeHex.slice(12, 15)}-${clockSeq}-${node}`
}

function copy(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped lang="scss">
.uuid-generator {
  .options {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    flex-wrap: wrap;

    .option-item {
      display: flex;
      align-items: center;
      gap: 12px;

      .label {
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        flex-shrink: 0;
      }
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .result {
    max-height: 400px;
    overflow-y: auto;

    .uuid-item {
      display: flex;
      align-items: center;
      padding: 8px 12px;
      background: #f8f9fa;
      border-radius: 6px;
      margin-bottom: 4px;

      .uuid-text {
        font-family: 'Courier New', monospace;
        font-size: 13px;
        color: #303133;
        flex: 1;
        word-break: break-all;
      }
    }
  }
}
</style>
