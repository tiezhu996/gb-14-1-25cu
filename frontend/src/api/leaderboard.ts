// 排行榜 API
import request from '../utils/request'
import type { LeaderboardEntry } from '../types'

export function getLeaderboard(period: string, limit = 50) {
  return request.get<unknown, LeaderboardEntry[]>('/leaderboard', { params: { period, limit } })
}
