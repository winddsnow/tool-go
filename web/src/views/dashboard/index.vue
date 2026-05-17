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
      <template #header>欢迎使用管理系统</template>
      <p>这是一个基于 GoFrame + Vue3 构建的企业级后台管理系统。</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { dashboardApi, DashboardStatsRes } from '@/api/dashboard'

const stats = ref<DashboardStatsRes>({
  user_count: 0,
  role_count: 0,
  online_user: 0,
  api_request: 0,
  total_visits: 0,
  user_visits: [],
})

const fetchStats = async () => {
  try {
    stats.value = await dashboardApi.getStats()
  } catch {
    // ignore
  }
}

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
