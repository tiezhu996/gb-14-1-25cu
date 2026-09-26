// 学习仪表盘 API
import request from '../utils/request'
import type { Dashboard } from '../types'

export function getDashboard() {
  return request.get<unknown, Dashboard>('/dashboard/me')
}
