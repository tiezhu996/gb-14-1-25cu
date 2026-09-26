// 统一请求封装：JWT 注入、错误码拦截、业务错误提示、统一解包
import axios, { AxiosError } from 'axios'
import type { ApiResponse } from '../types'

export interface ApiErrorPayload {
  code: number
  message: string
}

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('codelearn_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    // 统一解包后端响应 { code, message, data }，code=0 时返回 data
    const body = response.data as ApiResponse<unknown> | undefined
    if (body && body.code === 0) {
      return body.data as never
    }
    const err = new Error(body?.message || '请求失败') as Error & { code?: number; status?: number }
    err.code = body?.code
    err.status = response.status
    throw err
  },
  (error: AxiosError<ApiErrorPayload>) => {
    const payload = error.response?.data
    const message = payload?.message || error.message || '网络请求失败'
    const err = new Error(message) as Error & { code?: number; status?: number }
    err.code = payload?.code
    err.status = error.response?.status
    if (error.response?.status === 401) {
      localStorage.removeItem('codelearn_token')
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    throw err
  },
)

export default request
