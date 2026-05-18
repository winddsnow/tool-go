<template>
  <div class="qr-code-generator">
    <div class="input-area">
      <span class="label">文本或 URL</span>
      <el-input v-model="text" placeholder="请输入文本或 URL" size="large" @keyup.enter="generate" />
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="generate">生成二维码</el-button>
    </div>
    <div class="qr-result" v-if="showQR">
      <canvas ref="canvasRef"></canvas>
      <div class="qr-actions">
        <el-button type="success" size="small" @click="saveImage">保存为图片</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'

const text = ref('')
const showQR = ref(false)
const canvasRef = ref<HTMLCanvasElement | null>(null)

async function generate() {
  if (!text.value) {
    ElMessage.warning('请输入文本或 URL')
    return
  }
  showQR.value = true
  await nextTick()
  if (canvasRef.value) {
    try {
      await QRCode.toCanvas(canvasRef.value, text.value, { width: 256, margin: 2 })
    } catch {
      ElMessage.error('生成二维码失败')
    }
  }
}

function saveImage() {
  if (!canvasRef.value) return
  const link = document.createElement('a')
  link.download = 'qrcode.png'
  link.href = canvasRef.value.toDataURL('image/png')
  link.click()
}
</script>

<style scoped lang="scss">
.qr-code-generator {
  .input-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .qr-result {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;

    canvas {
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }
}
</style>
