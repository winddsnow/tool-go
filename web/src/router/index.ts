// index.ts — Vue Router 路由配置
// Vue Router 是 Vue.js 官方的路由管理库，负责 URL 与页面组件的映射
// 当用户点击链接或直接在地址栏输入 URL 时，路由会根据配置切换到对应的页面组件

// createRouter：创建路由实例
// createWebHistory：使用 HTML5 History 模式（URL 中不带 # 号，如 /user 而不是 /#/user）
// RouteRecordRaw：路由配置的类型定义
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
// useUserStore：用户状态 store，用于判断登录状态和角色权限
import { useUserStore } from '@/store/modules/user'
// authApi：调用后端接口获取用户信息
import { authApi } from '@/api/auth'

// routes 数组定义了所有路由规则
// 每个路由对象是一个 RouteRecordRaw，包含 path/component/meta 等属性
const routes: RouteRecordRaw[] = [
  // 登录页：不需要认证，所有用户都可以访问
  {
    path: '/login',
    name: 'Login',
    // 懒加载：() => import('...') 是动态导入语法
    // 只有当用户访问 /login 时才会加载 login/index.vue 的代码
    // 好处是首屏加载更快，不需要一次性下载所有页面的代码
    component: () => import('@/views/login/index.vue'),
    // meta：路由元信息，可以自由添加自定义数据
    // title：页面标题，会在导航守卫中设置到 document.title
    // requiresAuth: false 表示该页面不需要登录即可访问
    meta: { title: '登录', requiresAuth: false },
  },
  // 403 无权限页面：用户访问了没有权限的页面时跳转到此
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权限', requiresAuth: false },
  },
  // 主布局路由：包含侧边栏和顶栏等公共布局，子路由显示在布局的内容区
  {
    path: '/',
    // default.vue 是布局组件，内部包含 <router-view /> 用于渲染子路由
    component: () => import('@/layouts/default.vue'),
    // 访问 / 时自动重定向到 /tools
    redirect: '/tools',
    // children：嵌套路由，子路由的 path 会拼上父级的 path
    children: [
      {
        path: 'tools',          // 完整路径是 /tools
        name: 'Tools',
        component: () => import('@/views/tools/index.vue'),
        // icon：侧边栏菜单图标名，对应注册的 Element Plus 图标组件名
        meta: { title: '工具箱', icon: 'Tool' },
      },
      {
        path: 'user',           // 完整路径是 /user
        name: 'User',
        component: () => import('@/views/user/index.vue'),
        // requiresAuth: true 表示需要登录才能访问
        // roles 数组指定允许访问的角色，只有 super_admin 和 admin 可以查看用户管理页面
        meta: { title: '用户管理', icon: 'User', requiresAuth: true, roles: ['super_admin', 'admin'] },
      },
      {
        path: 'role',           // 完整路径是 /role
        name: 'Role',
        component: () => import('@/views/role/index.vue'),
        meta: { title: '角色管理', icon: 'Avatar', requiresAuth: true, roles: ['super_admin', 'admin'] },
      },
    ],
  },
  // 通配路由：匹配所有未定义路径（如 /xxx），重定向到 403 页面
  // :pathMatch(.*)* 会捕获所有路径
  {
    path: '/:pathMatch(.*)*',
    redirect: '/403',
  },
]

// 创建路由实例
const router = createRouter({
  // createWebHistory() 使用 HTML5 History 模式（推荐）
  // 效果：URL 为 /user 而不是 /#/user，更加美观
  // 注意：生产环境需要服务端配置，确保所有路径都返回 index.html（否则刷新会 404）
  history: createWebHistory(),
  routes,
})

// whiteList：白名单路由列表
// 未登录用户只能访问这些页面，访问其他页面会被重定向到登录页
const whiteList = ['/login', '/403', '/', '/tools']

// router.beforeEach：全局前置导航守卫
// 在每次路由切换之前执行，用于处理：页面标题设置、登录检查、角色权限检查
// to：即将进入的目标路由对象
// _from：当前正要离开的路由（下划线前缀表示该参数未使用）
// next：一个函数，调用 next() 放行，调用 next('/login') 重定向
router.beforeEach(async (to, _from, next) => {
  // 设置浏览器标签栏的标题
  document.title = `${to.meta.title || ''} - 瓦特的工具站`

  // 获取用户状态 store
  const userStore = useUserStore()
  // 判断用户是否有 token（是否已登录）
  const hasToken = userStore.token

  if (hasToken) {
    // ----- 用户已登录的情况 -----
    if (to.path === '/login') {
      // 如果已登录的用户访问登录页，直接跳转到首页
      next({ path: '/' })
    } else {
      // 访问其他页面
      if (userStore.roles.length === 0) {
        // roles 为空说明用户信息尚未加载（可能是页面刷新后 store 重置了）
        try {
          // 调用后端接口获取用户信息（包括角色列表）
          const userInfo = await authApi.getUserInfo()
          // 将用户信息存入 Pinia store，后续组件可以直接使用
          userStore.setUserInfo({
            userId: userInfo.user_id,
            username: userInfo.username,
            nickname: userInfo.nickname,
            roles: userInfo.roles,
          })
          // 重新执行当前导航，让后续的权限检查能读到刚设置的 roles
          next({ ...to, replace: true })
        } catch {
          // 获取用户信息失败（例如 token 过期），清除登录状态，跳回登录页
          userStore.logout()
          next(`/login?redirect=${to.path}`)
        }
      } else {
        // 用户信息已存在，检查角色权限
        // 从路由 meta 中读取需要的角色列表
        const requiredRoles = to.meta.roles as string[] | undefined
        if (requiredRoles && requiredRoles.length > 0) {
          // 当前页面设置了角色限制，检查用户是否拥有其中至少一个角色
          if (userStore.hasAnyRole(requiredRoles)) {
            next()  // 有权限，放行
          } else {
            next('/403')  // 无权限，跳转到 403 页面
          }
        } else {
          // 页面没有角色限制，直接放行
          next()
        }
      }
    }
  } else {
    // ----- 用户未登录的情况 -----
    if (whiteList.includes(to.path)) {
      // 访问的是白名单中的页面（登录页、403 页等），放行
      next()
    } else {
      // 访问受保护的页面，重定向到登录页
      // redirect 参数让登录成功后跳回原来想去的页面
      next(`/login?redirect=${to.path}`)
    }
  }
})

export default router
