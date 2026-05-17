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
        <el-menu-item index="/tools">
          <el-icon><Tool /></el-icon>
          <span>工具箱</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasAnyRole(['super_admin', 'admin'])" index="/user">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item v-if="userStore.hasAnyRole(['super_admin', 'admin'])" index="/role">
          <el-icon><Avatar /></el-icon>
          <span>角色管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-button v-if="userStore.hasAnyRole(['super_admin', 'admin'])" type="primary" link @click="router.push('/user')">
            <el-icon><Setting /></el-icon>
            管理后台
          </el-button>
          <span class="username">{{ userStore.nickname || userStore.username || '用户' }}</span>
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
  justify-content: space-between;
  padding: 0 20px;
  position: relative;
  z-index: 10;

  .header-left {
    display: flex;
    align-items: center;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .username {
      font-size: 14px;
      color: #606266;
    }
  }
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
