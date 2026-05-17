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
        @select="onMenuSelect"
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
// ============================================================
// Vue 3 组合式 API（Composition API）
// <script setup> 是 Vue 3 的语法糖，让代码更简洁。
// 顶层的 import / 变量 / 函数 自动暴露给模板使用。
// ============================================================

// ref     ：创建"响应式"数据（数字、字符串、布尔等），
//           通过 .value 读取/修改，模板中会自动展开（不用写 .value）
// computed：根据其他响应式数据自动计算新值，依赖变化时自动更新
import { ref, computed } from 'vue'

// useRoute()  -> 获取当前路由信息（路径、参数、meta 等）
// useRouter() -> 路由实例，用于编程式导航（push、replace 等）
import { useRoute, useRouter } from 'vue-router'

// Element Plus 图标组件，用于侧边栏的汉堡按钮
import { Expand } from '@element-plus/icons-vue'

// Pinia 状态管理（类似 Vuex）—— 全局用户状态（token、角色、昵称等）
import { useUserStore } from '@/store/modules/user'

// 登录/登出的 API 请求封装
import { authApi } from '@/api/auth'

// ----------------------------------------------------------
// route ：当前路由对象，可读取 route.path、route.meta.title 等
// router：路由实例，用 router.push('/login') 跳转页面
// userStore：Pinia store，存储用户登录状态、角色信息
// ----------------------------------------------------------
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// sidebarOpen：侧边栏是否展开（响应式布尔值）
// window.innerWidth >= 768：桌面端（≥768px）默认展开侧边栏，
// 移动端（<768px）默认收起，依赖汉堡按钮手动打开
const sidebarOpen = ref(window.innerWidth >= 768)

// isLoggedIn：计算属性，判断用户是否已登录
// computed 会根据依赖（userStore.token）自动重新计算
// !! 是 JavaScript 的布尔转换，将值转为 true/false
const isLoggedIn = computed(() => !!userStore.token)

// ----------------------------------------------------------
// onMenuSelect：点击侧边栏菜单项时触发
// @select="onMenuSelect" 是 el-menu 的选中事件
// 移动端（<768px）选中菜单后自动收起侧边栏
// ----------------------------------------------------------
const onMenuSelect = () => {
  if (window.innerWidth < 768) {
    sidebarOpen.value = false
  }
}

// ----------------------------------------------------------
// handleLogout：退出登录
// async/await 用于处理异步操作（如 API 请求）
// try/finally：无论 API 成功或失败，finally 块总会执行
// ----------------------------------------------------------
const handleLogout = async () => {
  try {
    // 调用后端退出接口（使 JWT 失效）
    await authApi.logout()
  } finally {
    // 清除 Pinia store 中的用户信息（token、角色等）
    userStore.logout()
    // 跳转到登录页面
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
