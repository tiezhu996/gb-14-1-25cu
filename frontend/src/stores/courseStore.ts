// 课程状态
import { create } from 'zustand'
import type { Chapter, Course } from '../types'
import * as courseApi from '../api/course'

interface CourseState {
  courses: Course[]
  total: number
  loading: boolean
  fetchCourses: (params?: { page?: number; page_size?: number; difficulty?: string; status?: string }) => Promise<void>
  createCourse: (payload: Partial<Course> & { chapters?: Chapter[] }) => Promise<Course>
}

export const useCourseStore = create<CourseState>((set) => ({
  courses: [],
  total: 0,
  loading: false,
  fetchCourses: async (params) => {
    set({ loading: true })
    try {
      const data = await courseApi.listCourses(params || {})
      set({ courses: data.list, total: data.total })
    } finally {
      set({ loading: false })
    }
  },
  createCourse: async (payload) => {
    const course = await courseApi.createCourse(payload)
    set((s) => ({ courses: [course, ...s.courses] }))
    return course
  },
}))
