<template>
  <div class="regex-tester">
    <div class="input-area">
      <span class="label">正则表达式</span>
      <el-input v-model="pattern" placeholder="/ 输入正则表达式 /" size="large" @input="testRegex" />
    </div>
    <div class="test-text-area">
      <span class="label">测试文本</span>
      <el-input v-model="testText" type="textarea" :rows="6" placeholder="请输入测试文本" size="large" @input="testRegex" />
    </div>
    <div class="options">
      <el-switch v-model="globalFlag" active-text="全局匹配 (g)" @change="testRegex" />
      <el-switch v-model="ignoreCaseFlag" active-text="忽略大小写 (i)" @change="testRegex" />
      <el-switch v-model="multilineFlag" active-text="多行模式 (m)" @change="testRegex" />
    </div>
    <div class="error-msg" v-if="error">{{ error }}</div>
    <div class="result" v-if="!error && matches.length > 0">
      <div class="match-count">匹配 {{ matches.length }} 处</div>
      <div class="match-list">
        <div v-for="(m, i) in matches" :key="i" class="match-item">
          <span class="match-index">#{{ i + 1 }}</span>
          <span class="match-text">{{ m[0] }}</span>
          <span class="match-pos">位置: {{ m.index }}</span>
        </div>
      </div>
    </div>
    <div class="no-match" v-else-if="!error && tested">无匹配结果</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const pattern = ref('')
const testText = ref('')
const globalFlag = ref(true)
const ignoreCaseFlag = ref(false)
const multilineFlag = ref(false)
const matches = ref<RegExpExecArray[]>([])
const error = ref('')
const tested = ref(false)

function testRegex() {
  error.value = ''
  matches.value = []
  tested.value = false
  if (!pattern.value || !testText.value) return
  tested.value = true
  try {
    const flags = `${globalFlag.value ? 'g' : ''}${ignoreCaseFlag.value ? 'i' : ''}${multilineFlag.value ? 'm' : ''}`
    const regex = new RegExp(pattern.value, flags)
    const results: RegExpExecArray[] = []
    let match: RegExpExecArray | null
    if (globalFlag.value) {
      while ((match = regex.exec(testText.value)) !== null) {
        results.push(match)
      }
    } else {
      match = regex.exec(testText.value)
      if (match) results.push(match)
    }
    matches.value = results
  } catch (e: any) {
    error.value = `正则错误: ${e.message}`
  }
}
</script>

<style scoped lang="scss">
.regex-tester {
  .input-area, .test-text-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }
  .options {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }
  .error-msg {
    color: #f56c6c;
    font-size: 13px;
    margin-bottom: 12px;
    padding: 8px 12px;
    background: #fef0f0;
    border-radius: 6px;
  }
  .result {
    .match-count {
      font-size: 14px;
      font-weight: 600;
      color: #409eff;
      margin-bottom: 12px;
    }
    .match-list {
      max-height: 300px;
      overflow-y: auto;
      .match-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 6px 12px;
        background: #f0f9ff;
        border-radius: 6px;
        margin-bottom: 4px;
        font-size: 13px;
        .match-index {
          color: #909399;
          min-width: 28px;
        }
        .match-text {
          font-family: 'Courier New', monospace;
          color: #303133;
          flex: 1;
          word-break: break-all;
        }
        .match-pos {
          color: #909399;
          font-size: 12px;
          flex-shrink: 0;
        }
      }
    }
  }
  .no-match {
    color: #909399;
    font-size: 14px;
    padding: 12px 0;
  }
}
</style>
