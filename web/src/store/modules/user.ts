import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userId = ref<number>(0)
  const username = ref<string>('')
  const nickname = ref<string>('')
  const roles = ref<string[]>([])

  function setToken(val: string) {
    token.value = val
    localStorage.setItem('token', val)
  }

  function setUserInfo(info: { userId: number; username: string; nickname: string; roles: string[] }) {
    userId.value = info.userId
    username.value = info.username
    nickname.value = info.nickname
    roles.value = info.roles
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
    localStorage.removeItem('token')
  }

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
