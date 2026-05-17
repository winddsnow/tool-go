<template>
  <div class="timestamp-converter">
    <div class="converter-section">
      <h3>时间戳 → 日期时间</h3>
      <div class="input-group">
        <el-input 
          v-model="timestampInput" 
          placeholder="请输入时间戳（秒或毫秒）" 
          size="large"
          @keyup.enter="convertToDateTime"
        >
          <template #append>
            <el-button @click="convertToDateTime">转换</el-button>
          </template>
        </el-input>
      </div>
      <div class="result-group" v-if="dateTimeResult">
        <span class="label">结果：</span>
        <span class="value">{{ dateTimeResult }}</span>
        <el-button size="small" link type="primary" @click="copyToClipboard(dateTimeResult)">复制</el-button>
      </div>
    </div>

    <el-divider />

    <div class="converter-section">
      <h3>日期时间 → 时间戳</h3>
      <div class="input-group">
        <el-date-picker
          v-model="dateTimeInput"
          type="datetime"
          placeholder="选择日期时间"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
          size="large"
          style="width: 100%"
        />
      </div>
      <div class="result-group" v-if="timestampResult">
        <span class="label">秒：</span>
        <span class="value">{{ timestampResult.seconds }}</span>
        <el-button size="small" link type="primary" @click="copyToClipboard(timestampResult.seconds)">复制</el-button>
      </div>
      <div class="result-group" v-if="timestampResult">
        <span class="label">毫秒：</span>
        <span class="value">{{ timestampResult.milliseconds }}</span>
        <el-button size="small" link type="primary" @click="copyToClipboard(timestampResult.milliseconds)">复制</el-button>
      </div>
    </div>

    <div class="current-time">
      <h3>当前时间</h3>
      <div class="time-display">{{ currentTime }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

const timestampInput = ref('')
const dateTimeInput = ref('')
const dateTimeResult = ref('')
const timestampResult = ref<{ seconds: string; milliseconds: string } | null>(null)
const currentTime = ref('')

let timer: number | null = null

const convertToDateTime = () => {
  if (!timestampInput.value) return
  
  let ts = Number(timestampInput.value)
  if (isNaN(ts)) {
    ElMessage.warning('请输入有效的时间戳')
    return
  }

  if (ts.toString().length === 10) {
    ts *= 1000
  }

  const date = new Date(ts)
  if (isNaN(date.getTime())) {
    ElMessage.warning('时间戳无效')
    return
  }

  const pad = (n: number) => n.toString().padStart(2, '0')
  const y = date.getFullYear()
  const m = pad(date.getMonth() + 1)
  const d = pad(date.getDate())
  const h = pad(date.getHours())
  const min = pad(date.getMinutes())
  const s = pad(date.getSeconds())

  dateTimeResult.value = `${y}-${m}-${d} ${h}:${min}:${s}`
}

const convertToTimestamp = () => {
  if (!dateTimeInput.value) {
    timestampResult.value = null
    return
  }

  const date = new Date(dateTimeInput.value)
  if (isNaN(date.getTime())) {
    ElMessage.warning('日期无效')
    return
  }

  timestampResult.value = {
    seconds: Math.floor(date.getTime() / 1000).toString(),
    milliseconds: date.getTime().toString(),
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const updateCurrentTime = () => {
  const now = new Date()
  const pad = (n: number) => n.toString().padStart(2, '0')
  currentTime.value = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`
}

onMounted(() => {
  updateCurrentTime()
  timer = window.setInterval(updateCurrentTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

import { watch } from 'vue'
watch(dateTimeInput, convertToTimestamp)
</script>

<style scoped lang="scss">
.timestamp-converter {
  padding: 20px;
  
  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
  }

  .converter-section {
    margin-bottom: 20px;
  }

  .result-group {
    margin-top: 12px;
    padding: 12px;
    background-color: #f5f7fa;
    border-radius: 8px;
    display: flex;
    align-items: center;
    gap: 8px;

    .label {
      color: #909399;
      font-size: 14px;
    }

    .value {
      color: #409eff;
      font-weight: 500;
      font-family: monospace;
      font-size: 16px;
    }
  }

  .current-time {
    margin-top: 30px;
    text-align: center;
    padding: 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    color: #fff;

    h3 {
      color: rgba(255, 255, 255, 0.8);
      margin-bottom: 10px;
    }

    .time-display {
      font-size: 28px;
      font-weight: bold;
      font-family: monospace;
      letter-spacing: 2px;
    }
  }
}
</style>
