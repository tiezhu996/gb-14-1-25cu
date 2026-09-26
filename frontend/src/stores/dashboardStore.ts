// 学习仪表盘状态
import { create } from 'zustand'
import type { Dashboard } from '../types'
import * as dashboardApi from '../api/dashboard'

interface DashboardState {
  data: Dashboard | null
  loading: boolean
  fetch: () => Promise<void>
}

export const useDashboardStore = create<DashboardState>((set) => ({
  data: null,
  loading: false,
  fetch: async () => {
    set({ loading: true })
    try {
      const data = await dashboardApi.getDashboard()
      set({ data })
    } finally {
      set({ loading: false })
    }
  },
}))
