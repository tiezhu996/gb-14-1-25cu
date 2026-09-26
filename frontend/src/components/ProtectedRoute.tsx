// 路由守卫：未登录跳转登录页；AdminRoute 校验管理员
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'

export function ProtectedRoute({ children }: { children: JSX.Element }) {
  const token = useAuthStore((s) => s.token)
  const location = useLocation()
  if (!token) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />
  }
  return children
}

export function AdminRoute({ children }: { children: JSX.Element }) {
  const token = useAuthStore((s) => s.token)
  const user = useAuthStore((s) => s.user)
  if (!token) {
    return <Navigate to="/login" replace />
  }
  if (user?.role !== 'admin') {
    return <Navigate to="/" replace />
  }
  return children
}
