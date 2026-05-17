<template>
  <div class="json-formatter">
    <div class="tool-section">
      <div class="input-area">
        <!--
          el-input 使用 type="textarea" 渲染为多行文本输入框
          v-model="input" 双向绑定输入内容
          :rows="8" 默认显示 8 行高度
        -->
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
        <!-- readonly 只读输出框，用户不能编辑，仅展示结果 -->
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
        <!-- :disabled="!output"：没有输出结果时禁用复制按钮 -->
        <el-button type="primary" link @click="copyResult" :disabled="!output">复制结果</el-button>
        <!-- v-if="errorMsg"：仅在出现错误时显示错误信息 -->
        <span v-if="errorMsg" class="error-msg">{{ errorMsg }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

// input：用户输入的原始 JSON 字符串
const input = ref('')
// output：格式化/压缩后的结果字符串
const output = ref('')
// errorMsg：错误提示信息，格式正确时为空字符串
const errorMsg = ref('')

// 格式化 JSON：将输入的 JSON 字符串解析后，用 2 空格缩进重新序列化
function formatJSON() {
  errorMsg.value = ''
  if (!input.value.trim()) {
    errorMsg.value = '请输入 JSON 字符串'
    return
  }
  try {
    // JSON.parse 将字符串解析为 JavaScript 对象（同时校验 JSON 格式是否正确）
    const parsed = JSON.parse(input.value)
    // JSON.stringify 的第三个参数为缩进空格数，2 表示每层缩进 2 个空格，实现美化效果
    output.value = JSON.stringify(parsed, null, 2)
  } catch (e: any) {
    // JSON.parse 失败说明输入不是合法的 JSON，捕获异常并提示错误信息
    errorMsg.value = 'JSON 格式错误: ' + e.message
    output.value = ''
  }
}

// 压缩 JSON：移除所有空格和换行，输出紧凑的单行 JSON
function compressJSON() {
  errorMsg.value = ''
  if (!input.value.trim()) {
    errorMsg.value = '请输入 JSON 字符串'
    return
  }
  try {
    const parsed = JSON.parse(input.value)
    // 不传缩进参数（或传 null），JSON.stringify 输出无空格的最小化字符串
    output.value = JSON.stringify(parsed)
  } catch (e: any) {
    errorMsg.value = 'JSON 格式错误: ' + e.message
    output.value = ''
  }
}

// 复制结果到剪贴板
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
