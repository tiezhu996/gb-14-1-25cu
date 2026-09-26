// 题目状态
import { create } from 'zustand'
import type { Problem, TestCase } from '../types'
import * as problemApi from '../api/problem'

interface ProblemState {
  problems: Problem[]
  total: number
  loading: boolean
  fetchProblems: (params?: { page?: number; page_size?: number; difficulty?: string; status?: string; tag?: string }) => Promise<void>
  createProblem: (payload: {
    title: string
    description: string
    difficulty: string
    languages: string[]
    tags?: string[]
    test_cases: TestCase[]
    time_limit?: number
    status?: string
  }) => Promise<Problem>
}

export const useProblemStore = create<ProblemState>((set) => ({
  problems: [],
  total: 0,
  loading: false,
  fetchProblems: async (params) => {
    set({ loading: true })
    try {
      const data = await problemApi.listProblems(params || {})
      set({ problems: data.list, total: data.total })
    } finally {
      set({ loading: false })
    }
  },
  createProblem: async (payload) => {
    const problem = await problemApi.createProblem(payload)
    set((s) => ({ problems: [problem, ...s.problems] }))
    return problem
  },
}))
