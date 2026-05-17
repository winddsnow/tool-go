import request from '@/utils/request'

export interface DashboardStatsRes {
  user_count: number
  role_count: number
  online_user: number
  api_request: number
}

export const dashboardApi = {
  getStats: () => request.get<DashboardStatsRes>('/dashboard/stats'),
}
