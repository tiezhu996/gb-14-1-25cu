// 认证状态：用户 + JWT（持久化到 localStorage）
import { create } from 'zustand'
import type { User } from '../types'
import { getMe } from '../api/user'

interface AuthState {
  user: User | null
  token: string
  setAuth: (token: string, user: User) => void
  setUser: (user: User) => void
  logout: () => void
  refreshMe: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem('codelearn_token') || '',
  setAuth: (token, user) => {
    localStorage.setItem('codelearn_token', token)
    set({ token, user })
  },
  setUser: (user) => set({ user }),
  logout: () => {
    localStorage.removeItem('codelearn_token')
    set({ user: null, token: '' })
  },
  refreshMe: async () => {
    try {
      const user = await getMe()
      set({ user })
    } catch {
      // 未登录时静默忽略
    }
  },
}))
