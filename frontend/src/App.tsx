// 应用路由：登录/注册 + 业务页面 + 管理后台（路由守卫）
import { useEffect } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Layout from './components/Layout'
import { ProtectedRoute, AdminRoute } from './components/ProtectedRoute'
import { useAuthStore } from './stores/authStore'
import Login from './pages/login/Login'
import Register from './pages/register/Register'
import Dashboard from './pages/dashboard/Dashboard'
import Courses from './pages/courses/Courses'
import CourseDetail from './pages/courses/CourseDetail'
import Problems from './pages/problems/Problems'
import ProblemDetail from './pages/problems/ProblemDetail'
import Submissions from './pages/submissions/Submissions'
import Leaderboard from './pages/leaderboard/Leaderboard'
import Achievements from './pages/achievements/Achievements'
import AdminCourses from './pages/admin/AdminCourses'
import AdminProblems from './pages/admin/AdminProblems'
import AdminUsers from './pages/admin/AdminUsers'
import AdminAudits from './pages/admin/AdminAudits'

export default function App() {
  useEffect(() => {
    // 已登录时刷新用户信息
    useAuthStore.getState().refreshMe()
  }, [])

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route path="/" element={<Dashboard />} />
          <Route path="/courses" element={<Courses />} />
          <Route path="/courses/:id" element={<CourseDetail />} />
          <Route path="/problems" element={<Problems />} />
          <Route path="/problems/:id" element={<ProblemDetail />} />
          <Route path="/submissions" element={<Submissions />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/achievements" element={<Achievements />} />
        </Route>
        <Route
          element={
            <AdminRoute>
              <Layout />
            </AdminRoute>
          }
        >
          <Route path="/admin/courses" element={<AdminCourses />} />
          <Route path="/admin/problems" element={<AdminProblems />} />
          <Route path="/admin/users" element={<AdminUsers />} />
          <Route path="/admin/audits" element={<AdminAudits />} />
        </Route>
        <Route path="*" element={<Login />} />
      </Routes>
    </BrowserRouter>
  )
}
