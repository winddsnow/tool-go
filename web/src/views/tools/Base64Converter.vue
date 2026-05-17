<template>
  <div class="base64-converter">
    <div class="tool-section">
      <!-- 方向选择：通过 radio 切换"编码"或"解码"模式 -->
      <div class="direction-select">
        <el-radio-group v-model="direction">
          <el-radio value="encode">编码 (Encode)</el-radio>
          <el-radio value="decode">解码 (Decode)</el-radio>
        </el-radio-group>
      </div>
      <div class="input-area">
        <!-- placeholder 随编码/解码模式自动切换 -->
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
        <!-- 上下互换：将输出结果交换到输入框，方便连续操作 -->
        <el-button @click="swapInputOutput">上下互换</el-button>
        <el-button @click="input = ''; output = ''">清空</el-button>
      </div>
      <div class="output-area">
        <!-- readonly 只读输出框，用户不可编辑 -->
        <el-input
          v-model="output"
          type="textarea"
          :rows="6"
          placeholder="输出结果"
          size="large"
          readonly
        />
      </div>
      <!-- v-if="output"：仅在有结果时才显示"复制结果"按钮 -->
      <div class="bottom-bar" v-if="output">
        <el-button type="primary" link @click="copyResult">复制结果</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

// direction：编码/解码方向，'encode' 为编码，'decode' 为解码
const direction = ref('encode')
const input = ref('')
const output = ref('')

// 计算属性：根据当前模式切换输入框的占位提示文字
const inputPlaceholder = computed(() =>
  direction.value === 'encode' ? '请输入要编码的文本' : '请输入 Base64 字符串'
)

// 执行编码或解码
function convert() {
  if (!input.value.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  try {
    if (direction.value === 'encode') {
      // 编码流程：
      // encodeURIComponent 将中文等非 ASCII 字符转为百分号编码（如 "你" → "%E4%BD%A0"）
      // unescape 将百分号编码解码为 Latin-1 字符串（每个字符 1 字节）
      // btoa 将 Latin-1 字符串转为 Base64
      output.value = btoa(unescape(encodeURIComponent(input.value)))
    } else {
      // 解码流程（反向操作）：
      // atob 将 Base64 解码为 Latin-1 字符串
      // escape 将 Latin-1 字符串转为百分号编码
      // decodeURIComponent 将百分号编码还原为原始 UTF-8 字符串（支持中文）
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

// 上下互换：将输出框的内容交换到输入框，同时输入内容清空到输出框
// 常用于先编码再继续修改的场景
function swapInputOutput() {
  if (output.value) {
    const tmp = input.value
    input.value = output.value
    output.value = tmp
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
