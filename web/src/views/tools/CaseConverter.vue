<template>
  <div class="case-converter">
    <div class="input-area">
      <span class="label">源文本</span>
      <el-input v-model="inputText" type="textarea" :rows="4" placeholder="请输入要转换的文本，如 hello world" size="large" />
    </div>
    <div class="actions" :class="{ 'is-mobile': isMobile }">
      <el-button v-for="conv in converters" :key="conv.key" @click="convert(conv.key)">
        <span class="conv-label">{{ conv.label }}</span>
        <span class="conv-example">{{ conv.example }}</span>
      </el-button>
    </div>
    <div class="output-area">
      <span class="label">转换结果</span>
      <el-input v-model="outputText" type="textarea" :rows="4" placeholder="点击上方按钮查看转换结果" size="large" readonly />
    </div>
    <el-button v-if="outputText" size="small" link type="primary" @click="copyResult">复制结果</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

const inputText = ref('')
const outputText = ref('')
const isMobile = ref(window.innerWidth < 768)

function onResize() { isMobile.value = window.innerWidth < 768 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

const converters = [
  { key: 'camel', label: 'camelCase', example: '→ helloWorld' },
  { key: 'pascal', label: 'PascalCase', example: '→ HelloWorld' },
  { key: 'snake', label: 'snake_case', example: '→ hello_world' },
  { key: 'kebab', label: 'kebab-case', example: '→ hello-world' },
  { key: 'constant', label: 'CONSTANT_CASE', example: '→ HELLO_WORLD' },
]

function splitWords(text: string): string[] {
  const trimmed = text.trim()
  if (!trimmed) return []
  const parts = trimmed.split(/[\s_-]+/)
  const words: string[] = []
  for (const part of parts) {
    if (!part) continue
    const matches = part.match(/[A-Z][a-z]+|[A-Z]+(?=[A-Z]|$|\d)|[a-z]+|\d+/g)
    if (matches) {
      words.push(...matches.map(w => w.toLowerCase()))
    } else {
      words.push(part.toLowerCase())
    }
  }
  return words.filter(Boolean)
}

function convert(type: string) {
  const words = splitWords(inputText.value)
  if (words.length === 0) {
    ElMessage.warning('请输入要转换的文本')
    return
  }
  switch (type) {
    case 'camel':
      outputText.value = words[0] + words.slice(1).map(w => w.charAt(0).toUpperCase() + w.slice(1)).join('')
      break
    case 'pascal':
      outputText.value = words.map(w => w.charAt(0).toUpperCase() + w.slice(1)).join('')
      break
    case 'snake':
      outputText.value = words.join('_')
      break
    case 'kebab':
      outputText.value = words.join('-')
      break
    case 'constant':
      outputText.value = words.map(w => w.toUpperCase()).join('_')
      break
  }
}

function copyResult() {
  navigator.clipboard.writeText(outputText.value).then(() => {
    ElMessage.success('已复制')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped lang="scss">
.case-converter {
  .input-area, .output-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }
  .actions {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
    flex-wrap: wrap;

    .el-button {
      display: flex;
      flex-direction: column;
      align-items: center;
      height: auto;
      padding: 8px 16px;

      .conv-label {
        font-size: 13px;
        font-weight: 600;
      }
      .conv-example {
        font-size: 11px;
        color: #909399;
        margin-top: 2px;
      }
    }

    &.is-mobile {
      .el-button {
        flex: 1;
        min-width: 80px;
      }
    }
  }
}
</style>
