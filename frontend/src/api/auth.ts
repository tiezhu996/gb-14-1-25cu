// 认证 API
import request from '../utils/request'
import type { LoginResponse } from '../types'

export function register(payload: { username: string; email: string; password: string; nickname?: string }) {
  return request.post<unknown, LoginResponse>('/auth/register', payload)
}

export function login(payload: { username: string; password: string }) {
  return request.post<unknown, LoginResponse>('/auth/login', payload)
}

export function signIn() {
  return request.post<unknown, { streak_days: number; first_today: boolean }>('/auth/sign-in')
}
