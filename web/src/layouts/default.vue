<template>
  <div class="layout-container">
    <el-aside width="220px">
      <div class="logo">管理系统</div>
      <el-menu
        :default-active="route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>工作台</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole('super_admin')" index="/user">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasRole('super_admin')" index="/role">
          <el-icon><Avatar /></el-icon>
          <span>角色管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-right">
          <span>{{ userStore.nickname || userStore.username || '管理员' }}</span>
          <el-button type="danger" link @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { authApi } from '@/api/auth'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const handleLogout = async () => {
  try {
    await authApi.logout()
  } finally {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped lang="scss">
.layout-container {
  height: 100%;
  display: flex;
}

.el-aside {
  background-color: #304156;
  color: #fff;

  .logo {
    height: 60px;
    line-height: 60px;
    text-align: center;
    font-size: 18px;
    font-weight: bold;
    color: #fff;
    border-bottom: 1px solid #1f2d3d;
  }
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
