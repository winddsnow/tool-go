// tools.ts — 工具箱相关的 API 接口
import request from '@/utils/request'

// MockDataReq：生成模拟数据时的请求参数
export interface MockDataReq {
  types: string[]   // 需要生成的字段类型数组，例如 ["name", "phone", "email"]
  count: number     // 生成的数据条数
}

// MockDataRes：后端返回的模拟数据
export interface MockDataRes {
  columns: string[]                   // 列名数组，例如 ["姓名", "手机号", "邮箱"]
  data: Record<string, string>[]      // 数据行数组，每行是一个对象，如 [{ "姓名": "张三", "手机号": "138..." }]
}

// toolsApi 对象管理工具箱相关功能
export const toolsApi = {
  // 生成模拟数据：POST /api/v1/tools/mock-data
  mockData: (data: MockDataReq) => request.post<MockDataRes>('/tools/mock-data', data),
}
