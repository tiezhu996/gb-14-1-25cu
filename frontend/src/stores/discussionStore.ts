// 讨论社区状态
import { create } from 'zustand'
import type { Discussion } from '../types'
import * as discussionApi from '../api/discussion'

interface DiscussionState {
  discussions: Discussion[]
  loading: boolean
  fetchDiscussions: (problemId: string, sort?: string) => Promise<void>
  createDiscussion: (problemId: string, payload: { title: string; content: string; code?: string; language?: string }) => Promise<void>
  vote: (id: string) => Promise<void>
  markBest: (problemId: string, id: string) => Promise<void>
}

export const useDiscussionStore = create<DiscussionState>((set, get) => ({
  discussions: [],
  loading: false,
  fetchDiscussions: async (problemId, sort) => {
    set({ loading: true })
    try {
      const data = await discussionApi.listDiscussions(problemId, sort)
      set({ discussions: data })
    } finally {
      set({ loading: false })
    }
  },
  createDiscussion: async (problemId, payload) => {
    const d = await discussionApi.createDiscussion(problemId, payload)
    set((s) => ({ discussions: [d, ...s.discussions] }))
  },
  vote: async (id) => {
    const d = await discussionApi.voteDiscussion(id)
    set((s) => ({ discussions: s.discussions.map((x) => (x.id === id ? d : x)) }))
  },
  markBest: async (problemId, id) => {
    await discussionApi.markBest(problemId, id)
    await get().fetchDiscussions(problemId, 'best')
  },
}))
