// 认证 Hook：读取本地用户、登录态判断、管理员判断
import { useEffect, useState } from 'react'
import { useAuthStore } from '../stores/authStore'

export function useAuth() {
  const { user, token, setUser, logout } = useAuthStore()
  const [ready, setReady] = useState(false)

  useEffect(() => {
    setReady(true)
  }, [])

  const isAdmin = user?.role === 'admin'
  return { user, token, isAdmin, setUser, logout, ready }
}
