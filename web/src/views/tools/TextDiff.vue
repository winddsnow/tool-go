<template>
  <div class="text-diff">
    <div class="diff-inputs" :class="{ 'is-mobile': isMobile }">
      <div class="input-area">
        <span class="label">原文本</span>
        <el-input v-model="oldText" type="textarea" :rows="8" placeholder="请输入原文本" size="large" />
      </div>
      <div class="input-area">
        <span class="label">新文本</span>
        <el-input v-model="newText" type="textarea" :rows="8" placeholder="请输入新文本" size="large" />
      </div>
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="computeDiff">对比</el-button>
      <el-button @click="clear">清空</el-button>
    </div>
    <div class="diff-output" v-if="diffResult.length > 0">
      <div
        v-for="(line, i) in diffResult"
        :key="i"
        class="diff-line"
        :class="line.type"
      >
        <span class="line-prefix">{{ line.prefix }}</span>
        <span class="line-text">{{ line.text }}</span>
      </div>
    </div>
    <div class="no-diff" v-else-if="compared">
      <el-empty :description="oldText === newText ? '两段文本完全相同' : '请输入文本进行对比'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const oldText = ref('')
const newText = ref('')
const diffResult = ref<{ type: 'add' | 'remove' | 'equal'; prefix: string; text: string }[]>([])
const compared = ref(false)
const isMobile = ref(window.innerWidth < 768)

function onResize() { isMobile.value = window.innerWidth < 768 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

function computeDiff() {
  compared.value = true
  if (!oldText.value && !newText.value) {
    diffResult.value = []
    return
  }
  const oldLines = oldText.value.split('\n')
  const newLines = newText.value.split('\n')
  const lcs = computeLCS(oldLines, newLines)
  diffResult.value = buildDiff(oldLines, newLines, lcs)
}

function computeLCS(a: string[], b: string[]): number[][] {
  const m = a.length, n = b.length
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (a[i - 1] === b[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
      }
    }
  }
  return dp
}

function buildDiff(a: string[], b: string[], dp: number[][]): { type: 'add' | 'remove' | 'equal'; prefix: string; text: string }[] {
  let i = a.length, j = b.length
  const temp: { type: 'add' | 'remove' | 'equal'; prefix: string; text: string }[] = []
  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && a[i - 1] === b[j - 1]) {
      temp.push({ type: 'equal', prefix: '  ', text: a[i - 1] })
      i--; j--
    } else if (j > 0 && (i === 0 || dp[i][j - 1] >= dp[i - 1][j])) {
      temp.push({ type: 'add', prefix: '+ ', text: b[j - 1] })
      j--
    } else {
      temp.push({ type: 'remove', prefix: '- ', text: a[i - 1] })
      i--
    }
  }
  return temp.reverse()
}

function clear() {
  oldText.value = ''
  newText.value = ''
  diffResult.value = []
  compared.value = false
}
</script>

<style scoped lang="scss">
.text-diff {
  .diff-inputs {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;

    &.is-mobile {
      flex-direction: column;
    }

    .input-area {
      flex: 1;
      .label {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 8px;
      }
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .diff-output {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 12px;
    max-height: 400px;
    overflow-y: auto;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.6;

    .diff-line {
      display: flex;
      padding: 2px 8px;
      border-radius: 4px;
      margin-bottom: 1px;

      &.add {
        background: #e6ffed;
        .line-prefix { color: #28a745; }
      }
      &.remove {
        background: #ffeef0;
        .line-prefix { color: #d73a49; }
      }
      &.equal {
        .line-prefix { color: #6a737d; }
      }

      .line-prefix {
        width: 24px;
        flex-shrink: 0;
      }
      .line-text {
        white-space: pre-wrap;
        word-break: break-all;
      }
    }
  }

  .no-diff {
    margin-top: 16px;
  }
}
</style>
