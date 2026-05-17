<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <template #header>用户总数</template>
          <div class="stat-value">{{ stats.user_count.toLocaleString() }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>角色数量</template>
          <div class="stat-value">{{ stats.role_count.toLocaleString() }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>总访问量</template>
          <div class="stat-value">{{ stats.total_visits.toLocaleString() }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>API请求</template>
          <div class="stat-value">{{ stats.api_request.toLocaleString() }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-card style="margin-top: 20px">
      <template #header>欢迎使用瓦特的工具站</template>
      <p>集合开发工具与后台管理为一体的工具箱。</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
// ============================================================
// Vue 3 组合式 API —— 本文件是仪表盘/首页
// ref       ：创建响应式数据（stats 是一个响应式对象）
// onMounted ：生命周期钩子，组件挂载到 DOM 后自动执行
// ============================================================
import { ref, onMounted } from 'vue'

// dashboardApi：仪表盘 API（获取统计数据）
// DashboardStatsRes：统计数据的 TypeScript 类型（包含各计数字段）
import { dashboardApi, DashboardStatsRes } from '@/api/dashboard'

// ----------------------------------------------------------
// stats：统计数据对象（响应式）
// 初始值全部为 0，组件挂载后从后端获取真实数据
// 模板中用 {{ stats.user_count.toLocaleString() }} 显示
// toLocaleString()：将数字转为带千分位分隔符的字符串（如 1,234）
// ----------------------------------------------------------
const stats = ref<DashboardStatsRes>({
  user_count: 0,    // 用户总数
  role_count: 0,    // 角色数量
  online_user: 0,   // 在线用户数
  api_request: 0,   // API 请求总数
  total_visits: 0,  // 总访问量
  user_visits: [],  // 按日期的用户访问趋势（暂时未使用）
})

// fetchStats：请求后端获取统计数据
const fetchStats = async () => {
  try {
    stats.value = await dashboardApi.getStats()
  } catch {
    // 请求失败时保持默认值（静默处理）
  }
}

// onMounted：组件挂载后立即请求数据
// 这里简写为 onMounted(fetchStats) —— 直接传入函数引用
// 相当于 onMounted(() => { fetchStats() })
onMounted(fetchStats)
</script>

<style scoped lang="scss">
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
  padding: 20px 0;
}
</style>
