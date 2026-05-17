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
      <!-- 算法选择：使用 el-radio-group 切换 md5 / sha1 / sha256 -->
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
      <!-- v-if="output"：仅在计算有结果时才显示输出区域 -->
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

// input：用户输入的待加密文本
const input = ref('')
// output：计算出的哈希值（十六进制字符串）
const output = ref('')
// algorithm：当前选中的哈希算法，默认 MD5
const algorithm = ref('md5')

// 计算属性：根据选中的算法返回中文名称
const algorithmLabel = computed(() => {
  const map: Record<string, string> = { md5: 'MD5', sha1: 'SHA1', sha256: 'SHA256' }
  return map[algorithm.value] || algorithm.value
})

// 计算哈希值的入口函数
async function computeHash() {
  if (!input.value.trim()) {
    ElMessage.warning('请输入要加密的文本')
    return
  }

  // TextEncoder 将字符串编码为 UTF-8 字节数组（Uint8Array），这是 Web Crypto API 需要的输入格式
  const encoder = new TextEncoder()
  const data = encoder.encode(input.value)

  let hashBuffer: ArrayBuffer
  // 根据选中的算法分支处理
  switch (algorithm.value) {
    case 'md5':
      // MD5 不在 Web Crypto API 标准中，使用纯 JS 手动实现
      hashBuffer = await computeMD5(data)
      break
    case 'sha1':
      // crypto.subtle.digest：浏览器内置的 Web Crypto API，硬件加速的哈希计算
      hashBuffer = await crypto.subtle.digest('SHA-1', data)
      break
    case 'sha256':
      hashBuffer = await crypto.subtle.digest('SHA-256', data)
      break
    default:
      return
  }

  // 将 ArrayBuffer（二进制数据）转换为十六进制字符串
  // new Uint8Array(hashBuffer)：将 ArrayBuffer 包装为可遍历的字节数组
  // Array.from() 将 Uint8Array 转为普通数组以便使用 .map()
  // 每个字节 b 用 b.toString(16) 转为十六进制，padStart(2, '0') 确保两位（如 "0f"）
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  output.value = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

// 纯 JavaScript 实现的 MD5 算法（浏览器 Web Crypto API 不提供 MD5）
async function computeMD5(data: Uint8Array): Promise<ArrayBuffer> {
  // MD5 核心逻辑中的四个辅助函数 F、G、H、I（位运算）
  const rotateLeft = (x: number, n: number) => (x << n) | (x >>> (32 - n))
  const F = (x: number, y: number, z: number) => (x & y) | (~x & z)
  const G = (x: number, y: number, z: number) => (x & z) | (y & ~z)
  const H = (x: number, y: number, z: number) => x ^ y ^ z
  const I = (x: number, y: number, z: number) => y ^ (x | ~z)

  // S 表：每轮循环的左移位数（每 16 轮切换一组）
  // T 表：基于正弦函数生成的 64 个常量
  const S: number[] = []
  const T: number[] = []
  for (let i = 0; i < 64; i++) {
    S[i] = [7, 12, 17, 22, 5, 9, 14, 20, 4, 11, 16, 23, 6, 10, 15, 21][
      Math.floor(i / 16) * 4 + (i % 4)
    ]
    T[i] = Math.floor(Math.abs(Math.sin(i + 1)) * 0x100000000)
  }

  // 数据填充与分块：MD5 按 512 位（64 字节）一块处理
  const blocks: number[][] = []
  const len = data.length
  const totalLen = len + 1 + 8 // 1 字节 0x80 填充 + 8 字节长度
  const blockCount = Math.ceil(totalLen / 64)
  for (let b = 0; b < blockCount; b++) {
    const block: number[] = []
    for (let i = 0; i < 16; i++) { // 每块 16 个 32 位整数
      let word = 0
      for (let j = 0; j < 4; j++) {
        const idx = b * 64 + i * 4 + j
        word |= (idx < len ? data[idx] : idx === len ? 0x80 : 0) << (j * 8)
      }
      // 最后一个块的倒数第二个 32 位存放消息长度（位为单位）的低 32 位
      if (b === blockCount - 1 && i === 14) word = (len * 8) & 0xffffffff
      // 最后一个块的最后一个 32 位存放消息长度的高 32 位
      if (b === blockCount - 1 && i === 15) word = Math.floor((len * 8) / 0x100000000)
      block.push(word)
    }
    blocks.push(block)
  }

  // MD5 初始向量（IV）
  let a = 0x67452301, b = 0xefcdab89, c = 0x98badcfe, d = 0x10325476

  // 主循环：处理每个 64 字节的数据块
  for (const block of blocks) {
    let aa = a, bb = b, cc = c, dd = d

    // 64 轮迭代
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

    // 将本块结果累加到初始向量上
    a = (a + aa) >>> 0
    b = (b + bb) >>> 0
    c = (c + cc) >>> 0
    d = (d + dd) >>> 0
  }

  // 将四个 32 位整数（a,b,c,d）按小端序拼接为 16 字节的 ArrayBuffer
  const result = new Uint8Array(16)
  for (let i = 0; i < 4; i++) {
    result[i] = (a >>> (i * 8)) & 0xff
    result[i + 4] = (b >>> (i * 8)) & 0xff
    result[i + 8] = (c >>> (i * 8)) & 0xff
    result[i + 12] = (d >>> (i * 8)) & 0xff
  }
  return result.buffer
}

// 复制哈希结果到剪贴板
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
