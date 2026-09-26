// 注册页
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { register } from '../../api/auth'
import { useAuthStore } from '../../stores/authStore'

export default function Register() {
  const [form, setForm] = useState({ username: '', email: '', password: '', nickname: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const setAuth = useAuthStore((s) => s.setAuth)
  const navigate = useNavigate()

  const update = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [k]: e.target.value }))

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const resp = await register(form)
      setAuth(resp.token, resp.user)
      navigate('/', { replace: true })
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-brand-600 to-indigo-800 p-6">
      <div className="w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl">
        <div className="text-center">
          <div className="text-4xl">🚀</div>
          <h1 className="mt-2 text-2xl font-bold text-gray-800">创建账号</h1>
        </div>
        <form onSubmit={submit} className="mt-6 space-y-4">
          {(
            [
              ['username', '用户名', '用户名（至少 3 个字符）'],
              ['email', '邮箱', 'you@example.com'],
              ['password', '密码', '密码（至少 6 位）'],
              ['nickname', '昵称', '选填'],
            ] as const
          ).map(([key, label, placeholder]) => (
            <div key={key}>
              <label className="mb-1 block text-sm font-medium text-gray-700">{label}</label>
              <input
                type={key === 'password' ? 'password' : key === 'email' ? 'email' : 'text'}
                value={form[key]}
                onChange={update(key)}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-brand-500 focus:outline-none"
                placeholder={placeholder}
                required={key !== 'nickname'}
              />
            </div>
          ))}
          {error && <div className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-600">{error}</div>}
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-brand-600 py-2.5 font-semibold text-white hover:bg-brand-700 disabled:opacity-60"
          >
            {loading ? '注册中...' : '注 册'}
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-gray-500">
          已有账号？
          <Link to="/login" className="font-medium text-brand-600 hover:underline">
            去登录
          </Link>
        </p>
      </div>
    </div>
  )
}
