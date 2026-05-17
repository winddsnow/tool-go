<template>
  <div class="tools-page">
    <div class="page-header">
      <h1>开发工具箱</h1>
      <p>高效、便捷的开发工具集合</p>
    </div>

    <div class="tools-grid">
      <div class="tool-card" @click="openTool('timestamp')">
        <div class="tool-icon">
          <el-icon :size="40"><Clock /></el-icon>
        </div>
        <div class="tool-info">
          <h3>时间戳转换</h3>
          <p>Unix 时间戳与日期时间双向转换</p>
        </div>
        <div class="tool-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>

      <div class="tool-card" @click="openTool('json')">
        <div class="tool-icon">
          <el-icon :size="40"><Edit /></el-icon>
        </div>
        <div class="tool-info">
          <h3>JSON 格式化</h3>
          <p>JSON 数据美化与校验</p>
        </div>
        <div class="tool-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>

      <div class="tool-card" @click="openTool('hash')">
        <div class="tool-icon">
          <el-icon :size="40"><Lock /></el-icon>
        </div>
        <div class="tool-info">
          <h3>哈希加密</h3>
          <p>MD5, SHA1, SHA256 加密工具</p>
        </div>
        <div class="tool-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>

      <div class="tool-card" @click="openTool('base64')">
        <div class="tool-icon">
          <el-icon :size="40"><Document /></el-icon>
        </div>
        <div class="tool-info">
          <h3>Base64 编解码</h3>
          <p>Base64 编码与解码工具</p>
        </div>
        <div class="tool-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>
    </div>

    <el-dialog v-model="toolVisible" :title="currentTool?.title" width="800px" destroy-on-close>
      <component :is="currentTool?.component" v-if="toolVisible" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, shallowRef } from 'vue'
import { Clock, Edit, Lock, Document, ArrowRight } from '@element-plus/icons-vue'
import TimestampConverter from './TimestampConverter.vue'
import JsonFormatter from './JsonFormatter.vue'
import HashEncryptor from './HashEncryptor.vue'
import Base64Converter from './Base64Converter.vue'

interface Tool {
  id: string
  title: string
  component: any
}

const tools: Record<string, Tool> = {
  timestamp: {
    id: 'timestamp',
    title: '时间戳转换工具',
    component: TimestampConverter,
  },
  json: {
    id: 'json',
    title: 'JSON 格式化工具',
    component: JsonFormatter,
  },
  hash: {
    id: 'hash',
    title: '哈希加密工具',
    component: HashEncryptor,
  },
  base64: {
    id: 'base64',
    title: 'Base64 编解码工具',
    component: Base64Converter,
  },
}

const toolVisible = ref(false)
const currentTool = ref<Tool | null>(null)

const openTool = (id: string) => {
  currentTool.value = tools[id] || null
  toolVisible.value = true
}
</script>

<style scoped lang="scss">
.tools-page {
  max-width: 1200px;
  margin: 0 auto;

  .page-header {
    text-align: center;
    margin-bottom: 40px;

    h1 {
      font-size: 32px;
      font-weight: 700;
      color: #1a1a2e;
      margin-bottom: 8px;
    }

    p {
      font-size: 16px;
      color: #8c8c8c;
    }
  }

  .tools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 24px;
  }

  .tool-card {
    background: #fff;
    border-radius: 16px;
    padding: 30px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    border: 1px solid #f0f0f0;
    position: relative;
    overflow: hidden;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      border-color: #409eff;

      .tool-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #fff;
      }

      .tool-arrow {
        opacity: 1;
        transform: translateX(0);
      }
    }

    &.disabled {
      cursor: not-allowed;
      opacity: 0.7;

      &:hover {
        transform: none;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        border-color: #f0f0f0;
      }
    }

    .tool-icon {
      width: 64px;
      height: 64px;
      border-radius: 12px;
      background: #f0f2f5;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #409eff;
      transition: all 0.3s ease;
    }

    .tool-info {
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
        margin: 0 0 4px;
      }

      p {
        font-size: 14px;
        color: #909399;
        margin: 0;
      }
    }

    .tool-arrow {
      position: absolute;
      top: 20px;
      right: 20px;
      opacity: 0;
      transform: translateX(-10px);
      transition: all 0.3s ease;
      color: #409eff;
    }

    .tool-badge {
      position: absolute;
      top: 16px;
      right: 16px;
      background: #f56c6c;
      color: #fff;
      font-size: 12px;
      padding: 2px 8px;
      border-radius: 10px;
    }
  }
}
</style>
