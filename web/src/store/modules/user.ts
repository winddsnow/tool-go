// user.ts — 用户状态管理（Pinia Store）
// Pinia 是 Vue 3 官方推荐的状态管理方案，类似于 Vuex
// Store（仓库）是一个全局共享的数据容器，所有组件都可以读取和修改里面的数据
// 相比组件内自己维护数据，Store 的好处是：不同页面之间共享数据无需层层传递

// defineStore 用于定义一个 Store
// useUserStore 是导出的函数，在组件中通过 useUserStore() 调用获取 store 实例
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { MenuTree } from '@/api/menu'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userId = ref<number>(0)
  const username = ref<string>('')
  const nickname = ref<string>('')
  const roles = ref<string[]>([])
  const menus = ref<MenuTree[]>(JSON.parse(localStorage.getItem('menus') || '[]'))

  function setToken(val: string) {
    token.value = val
    localStorage.setItem('token', val)
  }

  function setUserInfo(info: { userId: number; username: string; nickname: string; roles: string[]; menus?: MenuTree[] }) {
    userId.value = info.userId
    username.value = info.username
    nickname.value = info.nickname
    roles.value = info.roles
    menus.value = info.menus || []
    localStorage.setItem('menus', JSON.stringify(info.menus || []))
  }

  function hasRole(role: string): boolean {
    return roles.value.includes(role)
  }

  function hasAnyRole(roleList: string[]): boolean {
    return roleList.some(role => roles.value.includes(role))
  }

  function logout() {
    token.value = ''
    userId.value = 0
    username.value = ''
    nickname.value = ''
    roles.value = []
    menus.value = []
    localStorage.removeItem('token')
    localStorage.removeItem('menus')
  }

  return {
    token,
    userId,
    username,
    nickname,
    roles,
    menus,
    setToken,
    setUserInfo,
    hasRole,
    hasAnyRole,
    logout,
  }
})
