<template>
  <div class="base64-converter">
    <div class="tool-section">
      <div class="direction-select">
        <el-radio-group v-model="direction">
          <el-radio value="encode">编码 (Encode)</el-radio>
          <el-radio value="decode">解码 (Decode)</el-radio>
        </el-radio-group>
      </div>
      <div class="input-area">
        <el-input
          v-model="input"
          type="textarea"
          :rows="6"
          :placeholder="inputPlaceholder"
          size="large"
        />
      </div>
      <div class="action-bar">
        <el-button type="primary" @click="convert">执行{{ direction === 'encode' ? '编码' : '解码' }}</el-button>
        <el-button @click="swapInputOutput">上下互换</el-button>
        <el-button @click="input = ''; output = ''">清空</el-button>
      </div>
      <div class="output-area">
        <el-input
          v-model="output"
          type="textarea"
          :rows="6"
          placeholder="输出结果"
          size="large"
          readonly
        />
      </div>
      <div class="bottom-bar" v-if="output">
        <el-button type="primary" link @click="copyResult">复制结果</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const direction = ref('encode')
const input = ref('')
const output = ref('')

const inputPlaceholder = computed(() =>
  direction.value === 'encode' ? '请输入要编码的文本' : '请输入 Base64 字符串'
)

function convert() {
  if (!input.value.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  try {
    if (direction.value === 'encode') {
      output.value = btoa(unescape(encodeURIComponent(input.value)))
    } else {
      const decoded = decodeURIComponent(escape(atob(input.value.trim())))
      output.value = decoded
    }
  } catch {
    ElMessage.error(
      direction.value === 'encode' ? '编码失败，请检查输入内容' : '解码失败，请检查 Base64 字符串是否正确'
    )
    output.value = ''
  }
}

function swapInputOutput() {
  if (output.value) {
    const tmp = input.value
    input.value = output.value
    output.value = tmp
  }
}

function copyResult() {
  navigator.clipboard.writeText(output.value).then(() => {
    ElMessage.success('已复制到剪贴板')
  })
}
</script>

<style scoped>
.base64-converter {
  padding: 20px;
}

.direction-select {
  margin-bottom: 12px;
}

.action-bar {
  display: flex;
  gap: 12px;
  margin: 12px 0;
}

.bottom-bar {
  margin-top: 12px;
}
</style>
