// user.ts — 用户状态管理（Pinia Store）
// Pinia 是 Vue 3 官方推荐的状态管理方案，类似于 Vuex
// Store（仓库）是一个全局共享的数据容器，所有组件都可以读取和修改里面的数据
// 相比组件内自己维护数据，Store 的好处是：不同页面之间共享数据无需层层传递

// defineStore 用于定义一个 Store
// useUserStore 是导出的函数，在组件中通过 useUserStore() 调用获取 store 实例
import { defineStore } from 'pinia'
// ref 是 Vue 3 的响应式 API，用来声明可变的状态变量
// 当 ref 的值发生变化时，所有使用该变量的组件会自动重新渲染
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // ----- 响应式状态（ref）-----
  // ref<T>(初始值) 声明一个类型为 T 的响应式变量
  // .value 是 ref 的取值/赋值方式（在 template 中会自动解包，不需要写 .value）

  // token：JWT 认证令牌，从 localStorage 读取已有令牌（浏览器关闭后重新打开仍保留登录状态）
  // localStorage.getItem('token') 如果用户之前登录过，浏览器会存有 token，避免每次都要重新登录
  const token = ref<string>(localStorage.getItem('token') || '')
  // 当前登录用户的 ID
  const userId = ref<number>(0)
  // 当前登录用户的用户名
  const username = ref<string>('')
  // 当前登录用户的昵称（显示名称）
  const nickname = ref<string>('')
  // 当前登录用户的角色列表，例如 ["super_admin", "admin"]
  // 用于前端权限判断：哪些页面可以访问，哪些按钮可以显示
  const roles = ref<string[]>([])

  // ----- 操作方法 -----

  // setToken：保存登录令牌
  // 将 token 同时保存在内存（token.value）和 localStorage（持久化存储）中
  // 这样页面刷新后还能从 localStorage 恢复登录状态
  function setToken(val: string) {
    token.value = val
    localStorage.setItem('token', val)
  }

  // setUserInfo：设置用户信息（登录成功后由路由守卫调用）
  // 参数 info 包含用户的基本信息和角色列表
  function setUserInfo(info: { userId: number; username: string; nickname: string; roles: string[] }) {
    userId.value = info.userId
    username.value = info.username
    nickname.value = info.nickname
    roles.value = info.roles
  }

  // hasRole：判断当前用户是否拥有指定的单个角色
  // 例如 userStore.hasRole('super_admin') 返回 true 或 false
  function hasRole(role: string): boolean {
    return roles.value.includes(role)
  }

  // hasAnyRole：判断当前用户是否拥有指定角色列表中的任意一个
  // 例如 userStore.hasAnyRole(['super_admin', 'admin']) 只要满足其中一个就返回 true
  // .some() 是 JavaScript 数组方法：遍历数组，任一元素使回调返回 true 则整体返回 true
  function hasAnyRole(roleList: string[]): boolean {
    return roleList.some(role => roles.value.includes(role))
  }

  // logout：退出登录，清除所有用户状态
  // 同时删除 localStorage 中的 token，确保页面刷新后不会恢复已过期的登录
  function logout() {
    token.value = ''
    userId.value = 0
    username.value = ''
    nickname.value = ''
    roles.value = []
    localStorage.removeItem('token')
  }

  // 返回对象中暴露的内容可以被外部组件使用
  // 在组件中：const userStore = useUserStore() 后即可使用 userStore.token / userStore.hasRole() 等
  return {
    token,
    userId,
    username,
    nickname,
    roles,
    setToken,
    setUserInfo,
    hasRole,
    hasAnyRole,
    logout,
  }
})
