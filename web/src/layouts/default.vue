<template>
  <div class="layout-container">
    <div class="sidebar-overlay" v-if="sidebarOpen" @click="sidebarOpen = false" />

    <el-aside :width="sidebarOpen ? '220px' : '0'" class="sidebar">
      <div class="logo">瓦特的工具站</div>
      <el-menu
        :default-active="route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        @select="sidebarOpen = false"
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
          <el-button class="menu-btn" text @click="sidebarOpen = !sidebarOpen">
            <el-icon :size="22"><Expand /></el-icon>
          </el-button>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <template v-if="isLoggedIn">
            <el-button v-if="userStore.hasAnyRole(['super_admin', 'admin'])" type="primary" link @click="router.push('/user')">
              <el-icon><Setting /></el-icon>
              <span class="nav-text">管理后台</span>
            </el-button>
            <span class="username">{{ userStore.nickname || userStore.username }}</span>
            <el-button type="danger" link @click="handleLogout">退出</el-button>
          </template>
          <template v-else>
            <el-button type="primary" link @click="router.push('/login')">
              <el-icon><User /></el-icon>
              登录
            </el-button>
          </template>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <div class="layout-footer">
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">粤ICP备2025511523号</a>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Expand } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { authApi } from '@/api/auth'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const sidebarOpen = ref(window.innerWidth >= 768)
const isLoggedIn = computed(() => !!userStore.token)

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

.sidebar-overlay {
  display: none;

  @media (max-width: 768px) {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 99;
  }
}

.sidebar {
  background-color: #304156;
  color: #fff;
  overflow: hidden;
  transition: width 0.25s ease;
  position: relative;
  z-index: 100;

  @media (max-width: 768px) {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
  }

  .logo {
    height: 60px;
    line-height: 60px;
    text-align: center;
    font-size: 18px;
    font-weight: bold;
    color: #fff;
    border-bottom: 1px solid #1f2d3d;
    white-space: nowrap;
    overflow: hidden;
  }
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  position: relative;
  z-index: 10;

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .menu-btn {
    display: none;

    @media (max-width: 768px) {
      display: inline-flex;
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;

    .username {
      font-size: 14px;
      color: #606266;

      @media (max-width: 480px) {
        display: none;
      }
    }

    .nav-text {
      @media (max-width: 480px) {
        display: none;
      }
    }
  }
}

.el-main {
  background-color: #f0f2f5;
  padding: 16px;

  @media (max-width: 480px) {
    padding: 12px;
  }
}

.layout-footer {
  text-align: center;
  padding: 16px 20px;
  background: #fff;
  border-top: 1px solid #e6e6e6;
  font-size: 13px;

  a {
    color: #909399;
    text-decoration: none;

    &:hover {
      color: #409eff;
    }
  }
}
</style>
