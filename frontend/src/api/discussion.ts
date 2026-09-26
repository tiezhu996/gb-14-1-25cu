// 讨论社区 API
import request from '../utils/request'
import type { Discussion } from '../types'

export function listDiscussions(problemId: string, sort = 'best') {
  return request.get<unknown, Discussion[]>(`/problems/${problemId}/discussions`, { params: { sort } })
}

export function createDiscussion(problemId: string, payload: { title: string; content: string; code?: string; language?: string }) {
  return request.post<unknown, Discussion>(`/problems/${problemId}/discussions`, payload)
}

export function voteDiscussion(id: string) {
  return request.post<unknown, Discussion>(`/discussions/${id}/vote`)
}

export function markBest(problemId: string, id: string) {
  return request.put<unknown, unknown>(`/problems/${problemId}/discussions/${id}/best`)
}

export function hideDiscussion(id: string) {
  return request.put<unknown, unknown>(`/discussions/${id}/hide`)
}
