// 审计日志状态（管理员）
import { create } from 'zustand'
import type { AuditLog } from '../types'
import * as auditApi from '../api/audit'

interface AuditState {
  logs: AuditLog[]
  total: number
  loading: boolean
  fetch: (params?: { page?: number; page_size?: number; username?: string; entity?: string }) => Promise<void>
}

export const useAuditStore = create<AuditState>((set) => ({
  logs: [],
  total: 0,
  loading: false,
  fetch: async (params) => {
    set({ loading: true })
    try {
      const data = await auditApi.listAudits(params || {})
      set({ logs: data.list, total: data.total })
    } finally {
      set({ loading: false })
    }
  },
}))
