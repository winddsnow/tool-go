<template>
  <div class="hash-encryptor">
    <div class="tool-section">
      <div class="input-area">
        <el-input
          v-model="input"
          type="textarea"
          :rows="4"
          placeholder="请输入要加密的文本"
          size="large"
        />
      </div>
      <div class="algo-select">
        <span class="label">算法：</span>
        <el-radio-group v-model="algorithm">
          <el-radio value="md5">MD5</el-radio>
          <el-radio value="sha1">SHA1</el-radio>
          <el-radio value="sha256">SHA256</el-radio>
        </el-radio-group>
      </div>
      <div class="action-bar">
        <el-button type="primary" @click="computeHash">计算哈希</el-button>
        <el-button @click="input = ''; output = ''">清空</el-button>
      </div>
      <div class="output-area" v-if="output">
        <div class="result-row">
          <span class="label">{{ algorithmLabel }}：</span>
          <span class="value">{{ output }}</span>
          <el-button size="small" link type="primary" @click="copyResult">复制</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const input = ref('')
const output = ref('')
const algorithm = ref('md5')

const algorithmLabel = computed(() => {
  const map: Record<string, string> = { md5: 'MD5', sha1: 'SHA1', sha256: 'SHA256' }
  return map[algorithm.value] || algorithm.value
})

async function computeHash() {
  if (!input.value.trim()) {
    ElMessage.warning('请输入要加密的文本')
    return
  }

  const encoder = new TextEncoder()
  const data = encoder.encode(input.value)

  let hashBuffer: ArrayBuffer
  switch (algorithm.value) {
    case 'md5':
      hashBuffer = await computeMD5(data)
      break
    case 'sha1':
      hashBuffer = await crypto.subtle.digest('SHA-1', data)
      break
    case 'sha256':
      hashBuffer = await crypto.subtle.digest('SHA-256', data)
      break
    default:
      return
  }

  const hashArray = Array.from(new Uint8Array(hashBuffer))
  output.value = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

async function computeMD5(data: Uint8Array): Promise<ArrayBuffer> {
  const rotateLeft = (x: number, n: number) => (x << n) | (x >>> (32 - n))
  const F = (x: number, y: number, z: number) => (x & y) | (~x & z)
  const G = (x: number, y: number, z: number) => (x & z) | (y & ~z)
  const H = (x: number, y: number, z: number) => x ^ y ^ z
  const I = (x: number, y: number, z: number) => y ^ (x | ~z)

  const S: number[] = []
  const T: number[] = []
  for (let i = 0; i < 64; i++) {
    S[i] = [7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21][
      Math.floor(i / 16) * 4 + (i % 4)
    ]
    T[i] = Math.floor(Math.abs(Math.sin(i + 1)) * 0x100000000)
  }

  const blocks: number[][] = []
  const len = data.length
  const totalLen = len + 1 + 8
  const blockCount = Math.ceil(totalLen / 64)
  for (let b = 0; b < blockCount; b++) {
    const block: number[] = []
    for (let i = 0; i < 16; i++) {
      let word = 0
      for (let j = 0; j < 4; j++) {
        const idx = b * 64 + i * 4 + j
        word |= (idx < len ? data[idx] : idx === len ? 0x80 : 0) << (j * 8)
      }
      if (b === blockCount - 1 && i === 14) word = (len * 8) & 0xffffffff
      if (b === blockCount - 1 && i === 15) word = Math.floor((len * 8) / 0x100000000)
      block.push(word)
    }
    blocks.push(block)
  }

  let a = 0x67452301, b = 0xefcdab89, c = 0x98badcfe, d = 0x10325476

  for (const block of blocks) {
    let aa = a, bb = b, cc = c, dd = d

    for (let i = 0; i < 64; i++) {
      let f: number, g: number
      if (i < 16) { f = F(bb, cc, dd); g = i }
      else if (i < 32) { f = G(bb, cc, dd); g = (5 * i + 1) % 16 }
      else if (i < 48) { f = H(bb, cc, dd); g = (3 * i + 5) % 16 }
      else { f = I(bb, cc, dd); g = (7 * i) % 16 }

      const temp = dd
      dd = cc
      cc = bb
      bb = bb + rotateLeft(aa + f + T[i] + block[g], S[i])
      aa = temp
    }

    a = (a + aa) >>> 0
    b = (b + bb) >>> 0
    c = (c + cc) >>> 0
    d = (d + dd) >>> 0
  }

  const result = new Uint8Array(16)
  for (let i = 0; i < 4; i++) {
    result[i] = (a >>> (i * 8)) & 0xff
    result[i + 4] = (b >>> (i * 8)) & 0xff
    result[i + 8] = (c >>> (i * 8)) & 0xff
    result[i + 12] = (d >>> (i * 8)) & 0xff
  }
  return result.buffer
}

function copyResult() {
  navigator.clipboard.writeText(output.value).then(() => {
    ElMessage.success('已复制到剪贴板')
  })
}
</script>

<style scoped>
.hash-encryptor {
  padding: 20px;
}

.algo-select {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 12px 0;

  .label {
    color: #606266;
    font-size: 14px;
    white-space: nowrap;
  }
}

.action-bar {
  display: flex;
  gap: 12px;
  margin: 12px 0;
}

.result-row {
  margin-top: 12px;
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;

  .label {
    color: #909399;
    font-size: 14px;
    white-space: nowrap;
  }

  .value {
    color: #409eff;
    font-weight: 500;
    font-family: monospace;
    font-size: 14px;
    word-break: break-all;
    flex: 1;
  }
}
</style>
