// 课程 API
import request from '../utils/request'
import type { Chapter, Course, PageData } from '../types'

export function listCourses(params: { page?: number; page_size?: number; difficulty?: string; status?: string }) {
  return request.get<unknown, PageData<Course>>('/courses', { params })
}

export function getCourse(id: string) {
  return request.get<unknown, Course>(`/courses/${id}`)
}

export function createCourse(payload: Partial<Course> & { chapters?: Chapter[] }) {
  return request.post<unknown, Course>('/courses', payload)
}

export function updateCourse(id: string, payload: Partial<Course> & { chapters?: Chapter[] }) {
  return request.put<unknown, Course>(`/courses/${id}`, payload)
}

export function updateCourseStatus(id: string, status: string) {
  return request.put<unknown, unknown>(`/courses/${id}/status`, { status })
}

export function deleteCourse(id: string) {
  return request.delete<unknown, unknown>(`/courses/${id}`)
}

export function recordLearn(id: string, durationMinutes: number) {
  return request.post<unknown, unknown>(`/courses/${id}/learn`, { duration_minutes: durationMinutes })
}

export function completeCourse(id: string) {
  return request.post<unknown, unknown>(`/courses/${id}/complete`)
}
