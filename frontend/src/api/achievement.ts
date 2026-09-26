// 成就 API
import request from '../utils/request'
import type { Achievement } from '../types'

export function listAchievements() {
  return request.get<unknown, Achievement[]>('/achievements')
}

export function listMyAchievements() {
  return request.get<unknown, Achievement[]>('/achievements/me')
}
