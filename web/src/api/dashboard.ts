import request from '@/utils/request'

export interface UserVisitItem {
  username: string
  count: number
}

export interface DashboardStatsRes {
  user_count: number
  role_count: number
  online_user: number
  api_request: number
  total_visits: number
  user_visits: UserVisitItem[]
}

export const dashboardApi = {
  getStats: () => request.get<DashboardStatsRes>('/dashboard/stats'),
  trackPageView: (pagePath: string) => request.post('/pageview/track', { page_path: pagePath }),
  getPageViewStats: () => request.get<{ total_visits: number; user_visits: UserVisitItem[] }>('/pageview/stats'),
}
