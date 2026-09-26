// 审计日志 API（管理员）
import request from '../utils/request'
import type { AuditLog, PageData } from '../types'

export function listAudits(params: { page?: number; page_size?: number; username?: string; entity?: string }) {
  return request.get<unknown, PageData<AuditLog>>('/audits', { params })
}
