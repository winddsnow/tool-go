import { RouteRecordRaw } from 'vue-router'
import { MenuTree } from '@/api/menu'

const componentMap: Record<string, () => Promise<any>> = {
  'views/tools/index.vue': () => import('@/views/tools/index.vue'),
  'views/user/index.vue': () => import('@/views/user/index.vue'),
  'views/role/index.vue': () => import('@/views/role/index.vue'),
  'views/system/menu/index.vue': () => import('@/views/system/menu/index.vue'),
}

export function menuToRoutes(menus: MenuTree[]): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  for (const menu of menus) {
    if (menu.type === 3) continue
    if (menu.visible === 0) continue

    const component = componentMap[menu.component]
    if (!component) continue

    routes.push({
      path: menu.path.startsWith('/') ? menu.path.slice(1) : menu.path,
      name: menu.name,
      component,
      meta: {
        title: menu.name,
        icon: menu.icon,
        requiresAuth: true,
        menuId: menu.id,
      },
    })
  }
  return routes
}

export function flattenMenus(menus: MenuTree[]): MenuTree[] {
  const result: MenuTree[] = []
  for (const menu of menus) {
    result.push(menu)
    if (menu.children && menu.children.length > 0) {
      result.push(...flattenMenus(menu.children))
    }
  }
  return result
}
