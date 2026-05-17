// dashboard.ts — 仪表盘（首页概览）相关的 API 接口
import request from '@/utils/request'

// UserVisitItem：单个用户的访问量统计数据
export interface UserVisitItem {
  username: string   // 用户名
  count: number      // 该用户的访问次数
}

// DashboardStatsRes：仪表盘首页展示的统计数据
export interface DashboardStatsRes {
  user_count: number     // 系统用户总数
  role_count: number     // 系统角色总数
  online_user: number    // 当前在线用户数
  api_request: number    // API 请求总数（累计）
  total_visits: number   // 页面总访问量
  user_visits: UserVisitItem[]  // 各用户访问量排行列表
}

// dashboardApi 对象管理仪表盘相关的数据请求
export const dashboardApi = {
  // 获取仪表盘统计数据：GET /api/v1/dashboard/stats
  getStats: () => request.get<DashboardStatsRes>('/dashboard/stats'),
  // 记录页面访问：POST /api/v1/pageview/track
  // 当用户切换页面时调用，将访问路径发送给后端用于统计
  trackPageView: (pagePath: string) => request.post('/pageview/track', { page_path: pagePath }),
}
