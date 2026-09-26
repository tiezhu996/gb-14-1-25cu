// 成就状态
import { create } from 'zustand'
import type { Achievement } from '../types'
import * as achievementApi from '../api/achievement'

interface AchievementState {
  all: Achievement[]
  mine: Achievement[]
  loading: boolean
  fetchAll: () => Promise<void>
  fetchMine: () => Promise<void>
}

export const useAchievementStore = create<AchievementState>((set) => ({
  all: [],
  mine: [],
  loading: false,
  fetchAll: async () => {
    set({ loading: true })
    try {
      const all = await achievementApi.listAchievements()
      set({ all })
    } finally {
      set({ loading: false })
    }
  },
  fetchMine: async () => {
    const mine = await achievementApi.listMyAchievements()
    set({ mine })
  },
}))
