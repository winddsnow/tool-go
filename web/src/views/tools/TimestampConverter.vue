<template>
  <div class="timestamp-converter">
    <!-- 时区选择器：使用 el-select 下拉框切换不同时区 -->
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
        <!-- el-input 输入框绑定 timestampInput，回车键 (@keyup.enter) 触发转换 -->
        <el-input
          v-model="timestampInput"
          placeholder="请输入时间戳（秒或毫秒）"
          size="large"
          @keyup.enter="convertToDateTime"
        >
          <!-- #append 插槽：在输入框右侧追加"转换"按钮 -->
          <template #append>
            <el-button @click="convertToDateTime">转换</el-button>
          </template>
        </el-input>
      </div>
      <!-- v-if="dateTimeResult"：有转换结果才显示 -->
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
        <!--
          el-date-picker：Element Plus 日期时间选择器
          type="datetime" 同时选择日期和时间
          format 和 value-format 均设置为 "YYYY-MM-DD HH:mm:ss"
          :teleported="false"：下拉面板不追加到 body，避免在弹窗内出现 z-index 层级问题
          popper-class="mobile-dt-popper"：为手机端提供自定义样式修复
          @change 在值变化时自动触发转换
        -->
        <el-date-picker
          v-model="dateTimeInput"
          type="datetime"
          placeholder="选择日期时间"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
          size="large"
          style="width: 100%"
          :teleported="false"
          popper-class="mobile-dt-popper"
          @change="convertToTimestamp"
        />
      </div>
      <!-- 显示秒级时间戳 -->
      <div class="result-group" v-if="timestampResult">
        <span class="label">秒：</span>
        <span class="value">{{ timestampResult.seconds }}</span>
        <el-button size="small" link type="primary" @click="copyToClipboard(timestampResult.seconds)">复制</el-button>
      </div>
      <!-- 显示毫秒级时间戳 -->
      <div class="result-group" v-if="timestampResult">
        <span class="label">毫秒：</span>
        <span class="value">{{ timestampResult.milliseconds }}</span>
        <el-button size="small" link type="primary" @click="copyToClipboard(timestampResult.milliseconds)">复制</el-button>
      </div>
    </div>

    <!-- 当前时间实时时钟，背景使用了渐变色 -->
    <div class="current-time">
      <h3>当前时间（{{ timezoneLabel }}）</h3>
      <div class="time-display">{{ currentTime }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 导入 Vue 3 组合式 API
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'

// timezone 响应式变量，存储当前选中的时区
// Intl.DateTimeFormat().resolvedOptions().timeZone 获取浏览器默认时区（如 "Asia/Shanghai"）
const timezone = ref(Intl.DateTimeFormat().resolvedOptions().timeZone || 'Asia/Shanghai')

// 可选的时区列表，每个选项包含中文标签和 IANA 时区值
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

// 计算属性：从 timezones 数组中查找当前时区对应的中文标签
const timezoneLabel = computed(() => {
  const found = timezones.find(t => t.value === timezone.value)
  return found ? found.label : timezone.value
})

// 用户输入的时间戳字符串
const timestampInput = ref('')
// 用户选择的日期时间字符串
const dateTimeInput = ref('')
// 时间戳转日期后的结果字符串
const dateTimeResult = ref('')
// 日期转时间戳的结果，包含秒和毫秒两个值（可为 null 表示无结果）
const timestampResult = ref<{ seconds: string; milliseconds: string } | null>(null)
// 当前时间字符串，由定时器每秒更新
const currentTime = ref('')

// 定时器的 ID，用于组件卸载时清除定时器
let timer: number | null = null

// 核心工具函数：使用 Intl.DateTimeFormat 将 Date 对象按指定时区格式化为 "YYYY-MM-DD HH:mm:ss"
// formatToParts 返回一个数组（如 [{type:'year',value:'2024'}, ...]），方便自定义组装
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

  // 辅助函数：从 parts 数组中按 type 取出对应值，不足两位补 0
  const get = (type: string) => parts.find(p => p.type === type)?.value.padStart(2, '0') || '00'
  return `${get('year')}-${get('month')}-${get('day')} ${get('hour')}:${get('minute')}:${get('second')}`
}

// 将时间戳转换为日期时间
const convertToDateTime = () => {
  if (!timestampInput.value) return

  let ts = Number(timestampInput.value)
  if (isNaN(ts)) {
    ElMessage.warning('请输入有效的时间戳')
    return
  }

  // 时间戳位数判断：10 位数字是秒级时间戳，需要乘以 1000 转为毫秒
  // 13 位数字是毫秒级时间戳，直接使用
  if (ts.toString().length === 10) {
    ts *= 1000
  }

  const date = new Date(ts)
  if (isNaN(date.getTime())) {
    ElMessage.warning('时间戳无效')
    return
  }

  // 按当前选中的时区格式化结果
  dateTimeResult.value = formatInTimezone(date, timezone.value)
}

// 将日期时间转换为时间戳（秒和毫秒两个值）
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

  // getTime() 返回毫秒数，除以 1000 取整得到秒数
  timestampResult.value = {
    seconds: Math.floor(date.getTime() / 1000).toString(),
    milliseconds: date.getTime().toString(),
  }
}

// 复制到剪贴板的通用函数，使用浏览器原生 Clipboard API
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 更新当前实时时钟显示
const updateCurrentTime = () => {
  currentTime.value = formatInTimezone(new Date(), timezone.value)
}

// 监听 timezone 变化：切换时区时，重新计算已有结果和实时时钟
watch(timezone, () => {
  if (dateTimeResult.value) convertToDateTime()
  updateCurrentTime()
})

// 监听日期时间输入变化，自动转时间戳（无需手动点击）
watch(dateTimeInput, convertToTimestamp)

// 组件挂载后：立即显示当前时间，然后每秒刷新一次
onMounted(() => {
  updateCurrentTime()
  // setInterval 每隔 1000ms（1秒）更新当前时间，实现秒级走动的实时时钟
  timer = window.setInterval(updateCurrentTime, 1000)
})

// 组件卸载前：清除定时器，防止组件销毁后还在后台执行
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

:deep(.mobile-dt-popper) {
  @media (max-width: 480px) {
    width: 95vw !important;
    max-height: 70vh;
    overflow-y: auto;

    .el-time-spinner {
      width: 100%;
    }
    .el-date-picker__time-header {
      flex-wrap: wrap;
    }
  }
}
</style>
