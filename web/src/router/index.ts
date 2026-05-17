import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { authApi } from '@/api/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', requiresAuth: false },
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权限', requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/default.vue'),
    redirect: '/tools',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'tools',
        name: 'Tools',
        component: () => import('@/views/tools/index.vue'),
        meta: { title: '工具箱', icon: 'Tool' },
      },
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/user/index.vue'),
        meta: { title: '用户管理', icon: 'User', roles: ['super_admin', 'admin'] },
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/role/index.vue'),
        meta: { title: '角色管理', icon: 'Avatar', roles: ['super_admin', 'admin'] },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/403',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const whiteList = ['/login', '/403']

router.beforeEach(async (to, _from, next) => {
  document.title = `${to.meta.title || ''} - 管理系统`
  
  const userStore = useUserStore()
  const hasToken = userStore.token
  
  if (hasToken) {
    if (to.path === '/login') {
      next({ path: '/' })
    } else {
      if (userStore.roles.length === 0) {
        try {
          const userInfo = await authApi.getUserInfo()
          userStore.setUserInfo({
            userId: userInfo.user_id,
            username: userInfo.username,
            nickname: userInfo.nickname,
            roles: userInfo.roles,
          })
          next({ ...to, replace: true })
        } catch {
          userStore.logout()
          next(`/login?redirect=${to.path}`)
        }
      } else {
        const requiredRoles = to.meta.roles as string[] | undefined
        if (requiredRoles && requiredRoles.length > 0) {
          if (userStore.hasAnyRole(requiredRoles)) {
            next()
          } else {
            next('/403')
          }
        } else {
          next()
        }
      }
    }
  } else {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router
