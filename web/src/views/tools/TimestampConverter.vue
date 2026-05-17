<template>
  <div class="timestamp-converter">
    <div class="timezone-select">
      <span class="label">时区：</span>
      <el-select v-model="timezone" size="large" style="width: 240px">
        <el-option
          v-for="tz in timezones"
          :key="tz.value"
          :label="tz.label"
          :value="tz.value"
        />
      </el-select>
    </div>

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
      <h3>当前时间（{{ timezoneLabel }}）</h3>
      <div class="time-display">{{ currentTime }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'

const timezone = ref(Intl.DateTimeFormat().resolvedOptions().timeZone || 'Asia/Shanghai')

const timezones = [
  { label: 'UTC (协调世界时)', value: 'UTC' },
  { label: 'Asia/Shanghai (北京时间)', value: 'Asia/Shanghai' },
  { label: 'Asia/Tokyo (日本标准时间)', value: 'Asia/Tokyo' },
  { label: 'Asia/Seoul (韩国标准时间)', value: 'Asia/Seoul' },
  { label: 'Asia/Hong_Kong (香港时间)', value: 'Asia/Hong_Kong' },
  { label: 'Asia/Singapore (新加坡时间)', value: 'Asia/Singapore' },
  { label: 'America/New_York (美东时间)', value: 'America/New_York' },
  { label: 'America/Chicago (美中时间)', value: 'America/Chicago' },
  { label: 'America/Denver (美山区时间)', value: 'America/Denver' },
  { label: 'America/Los_Angeles (美西时间)', value: 'America/Los_Angeles' },
  { label: 'Europe/London (伦敦时间)', value: 'Europe/London' },
  { label: 'Europe/Berlin (柏林时间)', value: 'Europe/Berlin' },
  { label: 'Europe/Paris (巴黎时间)', value: 'Europe/Paris' },
  { label: 'Europe/Moscow (莫斯科时间)', value: 'Europe/Moscow' },
  { label: 'Australia/Sydney (悉尼时间)', value: 'Australia/Sydney' },
  { label: 'Pacific/Auckland (奥克兰时间)', value: 'Pacific/Auckland' },
]

const timezoneLabel = computed(() => {
  const found = timezones.find(t => t.value === timezone.value)
  return found ? found.label : timezone.value
})

const timestampInput = ref('')
const dateTimeInput = ref('')
const dateTimeResult = ref('')
const timestampResult = ref<{ seconds: string; milliseconds: string } | null>(null)
const currentTime = ref('')

let timer: number | null = null

function formatInTimezone(date: Date, tz: string): string {
  const parts = new Intl.DateTimeFormat('zh-CN', {
    timeZone: tz,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  }).formatToParts(date)

  const get = (type: string) => parts.find(p => p.type === type)?.value.padStart(2, '0') || '00'
  return `${get('year')}-${get('month')}-${get('day')} ${get('hour')}:${get('minute')}:${get('second')}`
}

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

  dateTimeResult.value = formatInTimezone(date, timezone.value)
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
  currentTime.value = formatInTimezone(new Date(), timezone.value)
}

watch(timezone, () => {
  if (dateTimeResult.value) convertToDateTime()
  updateCurrentTime()
})

watch(dateTimeInput, convertToTimestamp)

onMounted(() => {
  updateCurrentTime()
  timer = window.setInterval(updateCurrentTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped lang="scss">
.timestamp-converter {
  padding: 20px;

  .timezone-select {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 24px;

    .label {
      color: #606266;
      font-size: 14px;
      white-space: nowrap;
    }
  }

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
