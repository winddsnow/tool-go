<template>
  <div class="tools-page">
    <div class="stats-bar" v-if="stats">
      <el-row :gutter="16">
        <el-col :span="6">
          <div class="stat-item">
            <span class="stat-label">总访问量</span>
            <span class="stat-value">{{ stats.total_visits }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <span class="stat-label">用户总数</span>
            <span class="stat-value">{{ stats.user_count }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <span class="stat-label">角色数量</span>
            <span class="stat-value">{{ stats.role_count }}</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-item">
            <span class="stat-label">活跃访问</span>
            <span class="stat-value">{{ topUserLabel }}</span>
          </div>
        </el-col>
      </el-row>
    </div>

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

      <div class="tool-card" @click="openTool('password')">
        <div class="tool-icon">
          <el-icon :size="40"><Key /></el-icon>
        </div>
        <div class="tool-info">
          <h3>密码生成器</h3>
          <p>随机密码生成，支持自定义规则</p>
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
import { ref, computed, onMounted } from 'vue'
import { Clock, Edit, Lock, Document, Key, ArrowRight } from '@element-plus/icons-vue'
import { dashboardApi, DashboardStatsRes } from '@/api/dashboard'
import TimestampConverter from './TimestampConverter.vue'
import JsonFormatter from './JsonFormatter.vue'
import HashEncryptor from './HashEncryptor.vue'
import Base64Converter from './Base64Converter.vue'
import PasswordGenerator from './PasswordGenerator.vue'

interface Tool {
  id: string
  title: string
  component: any
}

const tools: Record<string, Tool> = {
  timestamp: { id: 'timestamp', title: '时间戳转换工具', component: TimestampConverter },
  json: { id: 'json', title: 'JSON 格式化工具', component: JsonFormatter },
  hash: { id: 'hash', title: '哈希加密工具', component: HashEncryptor },
  base64: { id: 'base64', title: 'Base64 编解码工具', component: Base64Converter },
  password: { id: 'password', title: '密码生成器', component: PasswordGenerator },
}

const toolVisible = ref(false)
const currentTool = ref<Tool | null>(null)
const stats = ref<DashboardStatsRes | null>(null)

const topUserLabel = computed(() => {
  if (!stats.value?.user_visits?.length) return '暂无'
  const top = stats.value.user_visits[0]
  return `${top.username} (${top.count})`
})

const openTool = (id: string) => {
  currentTool.value = tools[id] || null
  toolVisible.value = true
}

async function fetchStats() {
  try {
    stats.value = await dashboardApi.getStats()
  } catch {
    // ignore
  }
}

async function trackVisit() {
  try {
    await dashboardApi.trackPageView('/tools')
  } catch {
    // ignore
  }
}

onMounted(() => {
  trackVisit()
  fetchStats()
})
</script>

<style scoped lang="scss">
.tools-page {
  max-width: 1200px;
  margin: 0 auto;

  .stats-bar {
    background: #fff;
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 32px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    border: 1px solid #f0f0f0;

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;

      .stat-label {
        font-size: 13px;
        color: #909399;
      }

      .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: #409eff;
      }
    }
  }

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
  }
}
</style>
