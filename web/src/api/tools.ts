import request from '@/utils/request'

export interface MockDataReq {
  types: string[]
  count: number
}

export interface MockDataRes {
  columns: string[]
  data: Record<string, string>[]
}

export const toolsApi = {
  mockData: (data: MockDataReq) => request.post<MockDataRes>('/tools/mock-data', data),
}
