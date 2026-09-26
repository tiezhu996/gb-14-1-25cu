// 题目 API
import request from '../utils/request'
import type { PageData, Problem, TestCase } from '../types'

export function listProblems(params: { page?: number; page_size?: number; difficulty?: string; status?: string; tag?: string }) {
  return request.get<unknown, PageData<Problem>>('/problems', { params })
}

export function getProblem(id: string) {
  return request.get<unknown, Problem>(`/problems/${id}`)
}

export function createProblem(payload: {
  title: string
  description: string
  difficulty: string
  languages: string[]
  tags?: string[]
  test_cases: TestCase[]
  time_limit?: number
  status?: string
}) {
  return request.post<unknown, Problem>('/problems', payload)
}

export function updateProblem(id: string, payload: {
  title: string
  description: string
  difficulty: string
  languages: string[]
  tags?: string[]
  test_cases: TestCase[]
  time_limit?: number
  status?: string
}) {
  return request.put<unknown, Problem>(`/problems/${id}`, payload)
}

export function updateProblemStatus(id: string, status: string) {
  return request.put<unknown, unknown>(`/problems/${id}/status`, { status })
}

export function deleteProblem(id: string) {
  return request.delete<unknown, unknown>(`/problems/${id}`)
}
