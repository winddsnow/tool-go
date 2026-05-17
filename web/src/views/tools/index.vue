<template>
  <!-- 整个工具页面的根容器 -->
  <div class="tools-page">
    <!-- v-if="stats"：只有当 stats 数据加载成功（不为 null）时，才渲染统计栏 -->
    <div class="stats-bar" v-if="stats">
      <el-row :gutter="[16, 16]">
        <!-- :xs="12" 表示屏幕宽度 < 768px 时，每个项占 12 列（一行 2 个） -->
        <!-- :sm="6"  表示屏幕宽度 >= 768px 时，每个项占 6 列（一行 4 个） -->
        <!-- 这种写法实现了响应式布局，适配手机和桌面 -->
        <el-col :xs="12" :sm="6">
          <div class="stat-item">
            <span class="stat-label">总访问量</span>
            <span class="stat-value">{{ stats.total_visits }}</span>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-item">
            <span class="stat-label">用户总数</span>
            <span class="stat-value">{{ stats.user_count }}</span>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-item">
            <span class="stat-label">角色数量</span>
            <span class="stat-value">{{ stats.role_count }}</span>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
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

    <!-- 工具卡片网格，每个卡片点击后弹出对应的工具弹窗 -->
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

      <div class="tool-card" @click="openTool('mockdata')">
        <div class="tool-icon">
          <el-icon :size="40"><List /></el-icon>
        </div>
        <div class="tool-info">
          <h3>随机数据生成器</h3>
          <p>生成姓名/手机/身份证/护照等模拟数据</p>
        </div>
        <div class="tool-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>
    </div>

    <!--
      el-dialog：Element Plus 的弹窗组件
      v-model="toolVisible"：控制弹窗显示/隐藏
      :title="currentTool?.title"：动态显示当前工具的标题
      :width="isMobile ? '95%' : '800px'"：手机上宽度占 95%，桌面固定 800px
      destroy-on-close：关闭弹窗时销毁内部组件，确保每次打开都是全新实例
    -->
    <el-dialog v-model="toolVisible" :title="currentTool?.title" :width="isMobile ? '95%' : '800px'" destroy-on-close>
      <!--
        <component :is="...">：Vue 动态组件语法
        :is 绑定当前选中的工具组件（如 TimestampConverter、JsonFormatter 等）
        v-if="toolVisible"：仅在弹窗显示时才渲染组件，配合 destroy-on-close 避免残留
      -->
      <component :is="currentTool?.component" v-if="toolVisible" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
// Vue 3 组合式 API：ref 创建响应式数据，computed 创建计算属性
// onMounted / onUnmounted 是生命周期钩子，分别在组件挂载和卸载时触发
import { ref, computed, onMounted, onUnmounted } from 'vue'
// 从 Element Plus 图标库导入所需图标
import { Clock, Edit, Lock, Document, Key, List, ArrowRight } from '@element-plus/icons-vue'
// Pinia 状态管理：获取用户登录状态（token、角色等）
import { useUserStore } from '@/store/modules/user'
// 仪表盘 API：获取统计数据和记录页面访问
import { dashboardApi, DashboardStatsRes } from '@/api/dashboard'
// 导入 6 个子工具组件，用于动态渲染
import TimestampConverter from './TimestampConverter.vue'
import JsonFormatter from './JsonFormatter.vue'
import HashEncryptor from './HashEncryptor.vue'
import Base64Converter from './Base64Converter.vue'
import PasswordGenerator from './PasswordGenerator.vue'
import MockDataGenerator from './MockDataGenerator.vue'

// TypeScript 接口：定义工具对象的类型结构
interface Tool {
  id: string
  title: string
  component: any
}

// 工具注册表：用 Record 对象存储所有工具，key 是工具标识，value 包含标题和组件
const tools: Record<string, Tool> = {
  timestamp: { id: 'timestamp', title: '时间戳转换工具', component: TimestampConverter },
  json: { id: 'json', title: 'JSON 格式化工具', component: JsonFormatter },
  hash: { id: 'hash', title: '哈希加密工具', component: HashEncryptor },
  base64: { id: 'base64', title: 'Base64 编解码工具', component: Base64Converter },
  password: { id: 'password', title: '密码生成器', component: PasswordGenerator },
  mockdata: { id: 'mockdata', title: '随机数据生成器', component: MockDataGenerator },
}

// 弹窗显示/隐藏状态，初始为 false（隐藏）
const toolVisible = ref(false)
// 当前选中的工具对象，点击卡片时赋值
const currentTool = ref<Tool | null>(null)
// 仪表盘统计数据，初始为 null，加载成功后有值后才显示统计栏
const stats = ref<DashboardStatsRes | null>(null)
// 记录窗口宽度，用于计算是否为移动端
const windowWidth = ref(window.innerWidth)
// isMobile 计算属性：窗口宽度小于 768px 时认为当前是移动设备
const isMobile = computed(() => windowWidth.value < 768)
// 获取用户状态（token、角色等），用于判断是否登录
const userStore = useUserStore()

// 窗口尺寸变化时的回调函数，更新 windowWidth 的值
function onResize() { windowWidth.value = window.innerWidth }
// 组件挂载后：监听窗口 resize 事件
onMounted(() => window.addEventListener('resize', onResize))
// 组件卸载前：移除 resize 事件监听，防止内存泄漏
onUnmounted(() => window.removeEventListener('resize', onResize))

// 计算属性：获取访问量最高的用户及其访问次数，用于"活跃访问"统计项
const topUserLabel = computed(() => {
  if (!stats.value?.user_visits?.length) return '暂无'
  const top = stats.value.user_visits[0]
  return `${top.username} (${top.count})`
})

// 打开工具：点击卡片时调用，设置当前工具并显示弹窗
const openTool = (id: string) => {
  currentTool.value = tools[id] || null
  toolVisible.value = true
}

// 异步拉取统计面板数据（仅在用户已登录时执行）
async function fetchStats() {
  if (!userStore.token) return
  try {
    stats.value = await dashboardApi.getStats()
  } catch {
    // ignore
  }
}

// 记录页面访问（即使用户未登录也记录，用于统计匿名访客）
async function trackVisit() {
  try {
    await dashboardApi.trackPageView('/tools')
  } catch {
    // ignore
  }
}

// 组件挂载后自动执行：先记录访问，再加载统计数据
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

        @media (max-width: 480px) {
          font-size: 18px;
        }
      }
    }
  }

  .page-header {
    text-align: center;
    margin-bottom: 32px;

    h1 {
      font-size: 28px;
      font-weight: 700;
      color: #1a1a2e;
      margin-bottom: 8px;

      @media (max-width: 480px) {
        font-size: 22px;
      }
    }

    p {
      font-size: 16px;
      color: #8c8c8c;
    }
  }

  .tools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;

    @media (max-width: 480px) {
      grid-template-columns: 1fr;
      gap: 12px;
    }
  }

  .tool-card {
    background: #fff;
    border-radius: 16px;
    padding: 20px;
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
