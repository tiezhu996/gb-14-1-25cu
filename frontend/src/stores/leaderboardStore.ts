// 排行榜状态
import { create } from 'zustand'
import type { LeaderboardEntry } from '../types'
import * as leaderboardApi from '../api/leaderboard'

interface LeaderboardState {
  entries: LeaderboardEntry[]
  loading: boolean
  fetch: (period: string) => Promise<void>
}

export const useLeaderboardStore = create<LeaderboardState>((set) => ({
  entries: [],
  loading: false,
  fetch: async (period) => {
    set({ loading: true })
    try {
      const entries = await leaderboardApi.getLeaderboard(period)
      set({ entries })
    } finally {
      set({ loading: false })
    }
  },
}))
