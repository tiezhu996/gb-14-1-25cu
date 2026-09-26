// 提交评测 API
import request from '../utils/request'
import type { PageData, Submission } from '../types'

export function submitCode(problemId: string, payload: { language: string; code: string }) {
  return request.post<unknown, Submission>(`/problems/${problemId}/submit`, payload)
}

export function getSubmission(id: string) {
  return request.get<unknown, Submission>(`/submissions/${id}`)
}

export function listSubmissions(params: { page?: number; page_size?: number; problem_id?: string; status?: string }) {
  return request.get<unknown, PageData<Submission>>('/submissions', { params })
}
