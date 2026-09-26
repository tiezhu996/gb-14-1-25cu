// 登录页
import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { login } from '../../api/auth'
import { useAuthStore } from '../../stores/authStore'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const setAuth = useAuthStore((s) => s.setAuth)
  const navigate = useNavigate()
  const location = useLocation() as { state?: { from?: string } }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const resp = await login({ username, password })
      setAuth(resp.token, resp.user)
      navigate(location.state?.from || '/', { replace: true })
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const fillDemo = (u: string, p: string) => {
    setUsername(u)
    setPassword(p)
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-brand-600 to-indigo-800 p-6">
      <div className="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl">
        <div className="text-center">
          <div className="text-4xl">💻</div>
          <h1 className="mt-2 text-2xl font-bold text-gray-800">CodeLearn 在线编程学习平台</h1>
          <p className="mt-1 text-sm text-gray-500">登录后开始你的编程学习之旅</p>
        </div>
        <form onSubmit={submit} className="mt-6 space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">用户名</label>
            <input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-brand-500 focus:outline-none"
              placeholder="请输入用户名"
              required
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">密码</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-brand-500 focus:outline-none"
              placeholder="请输入密码"
              required
            />
          </div>
          {error && <div className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-600">{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-brand-600 py-2.5 font-semibold text-white hover:bg-brand-700 disabled:opacity-60"
          >
            {loading ? '登录中...' : '登 录'}
          </button>
        </form>
        <div className="mt-4 flex items-center justify-between text-sm">
          <button onClick={() => fillDemo('student', 'student123')} className="text-brand-600 hover:underline">
            学生演示账号
          </button>
          <button onClick={() => fillDemo('admin', 'admin123')} className="text-brand-600 hover:underline">
            管理员演示账号
          </button>
        </div>
        <p className="mt-4 text-center text-sm text-gray-500">
          还没有账号？
          <Link to="/register" className="font-medium text-brand-600 hover:underline">
            立即注册
          </Link>
        </p>
      </div>
    </div>
  )
}
