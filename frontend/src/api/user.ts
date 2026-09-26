// 用户 API
import request from '../utils/request'
import type { PageData, User } from '../types'

export function getMe() {
  return request.get<unknown, User>('/users/me')
}

export function updateProfile(payload: { nickname?: string; avatar?: string; email?: string }) {
  return request.put<unknown, User>('/users/me', payload)
}

export function listUsers(params: { page?: number; page_size?: number; role?: string; status?: string }) {
  return request.get<unknown, PageData<User>>('/users', { params })
}

export function updateUserStatus(id: string, status: string) {
  return request.put<unknown, unknown>(`/users/${id}/status`, { status })
}

export function updateUserRole(id: string, role: string) {
  return request.put<unknown, unknown>(`/users/${id}/role`, { role })
}
