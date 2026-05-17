<template>
  <div class="json-formatter">
    <div class="tool-section">
      <div class="input-area">
        <el-input
          v-model="input"
          type="textarea"
          :rows="8"
          placeholder="请输入 JSON 字符串"
          size="large"
        />
      </div>
      <div class="action-bar">
        <el-button type="primary" @click="formatJSON">格式化</el-button>
        <el-button @click="compressJSON">压缩</el-button>
        <el-button @click="input = ''">清空</el-button>
      </div>
      <div class="output-area">
        <el-input
          v-model="output"
          type="textarea"
          :rows="8"
          placeholder="输出结果"
          size="large"
          readonly
        />
      </div>
      <div class="bottom-bar">
        <el-button type="primary" link @click="copyResult" :disabled="!output">复制结果</el-button>
        <span v-if="errorMsg" class="error-msg">{{ errorMsg }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const input = ref('')
const output = ref('')
const errorMsg = ref('')

function formatJSON() {
  errorMsg.value = ''
  if (!input.value.trim()) {
    errorMsg.value = '请输入 JSON 字符串'
    return
  }
  try {
    const parsed = JSON.parse(input.value)
    output.value = JSON.stringify(parsed, null, 2)
  } catch (e: any) {
    errorMsg.value = 'JSON 格式错误: ' + e.message
    output.value = ''
  }
}

function compressJSON() {
  errorMsg.value = ''
  if (!input.value.trim()) {
    errorMsg.value = '请输入 JSON 字符串'
    return
  }
  try {
    const parsed = JSON.parse(input.value)
    output.value = JSON.stringify(parsed)
  } catch (e: any) {
    errorMsg.value = 'JSON 格式错误: ' + e.message
    output.value = ''
  }
}

function copyResult() {
  navigator.clipboard.writeText(output.value).then(() => {
    ElMessage.success('已复制到剪贴板')
  })
}
</script>

<style scoped>
.json-formatter {
  padding: 20px;
}

.action-bar {
  display: flex;
  gap: 12px;
  margin: 12px 0;
}

.bottom-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.error-msg {
  color: #f56c6c;
  font-size: 13px;
}
</style>
