<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <template #header>用户总数</template>
          <div class="stat-value">1,234</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>角色数量</template>
          <div class="stat-value">12</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>API请求</template>
          <div class="stat-value">56,789</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>在线用户</template>
          <div class="stat-value">89</div>
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
