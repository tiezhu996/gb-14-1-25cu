// 全局布局：侧边导航 + 顶栏（管理员可见管理菜单，按钮显隐由角色控制）
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

const navItems = [
  { to: '/', label: '学习仪表盘', icon: '📊', end: true },
  { to: '/courses', label: '课程中心', icon: '📚' },
  { to: '/problems', label: '题库练习', icon: '🧩' },
  { to: '/submissions', label: '我的提交', icon: '📝' },
  { to: '/leaderboard', label: '排行榜', icon: '🏆' },
  { to: '/achievements', label: '成就徽章', icon: '🎖️' },
]

const adminItems = [
  { to: '/admin/courses', label: '课程管理', icon: '⚙️' },
  { to: '/admin/problems', label: '题目管理', icon: '🛠️' },
  { to: '/admin/users', label: '用户管理', icon: '👥' },
  { to: '/admin/audits', label: '审计日志', icon: '🔍' },
]

export default function Layout() {
  const { user, isAdmin, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="flex min-h-screen">
      <aside className="fixed inset-y-0 left-0 w-60 bg-gray-900 text-gray-300">
        <div className="flex items-center gap-2 px-5 py-5">
          <span className="text-2xl">💻</span>
          <span className="text-lg font-bold text-white">CodeLearn</span>
        </div>
        <nav className="mt-2 space-y-1 px-3">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                `flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition ${
                  isActive ? 'bg-brand-600 text-white' : 'hover:bg-gray-800 hover:text-white'
                }`
              }
            >
              <span>{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
          {isAdmin && (
            <>
              <div className="px-3 pt-4 text-xs uppercase tracking-wider text-gray-500">管理后台</div>
              {adminItems.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  className={({ isActive }) =>
                    `flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition ${
                      isActive ? 'bg-brand-600 text-white' : 'hover:bg-gray-800 hover:text-white'
                    }`
                  }
                >
                  <span>{item.icon}</span>
                  {item.label}
                </NavLink>
              ))}
            </>
          )}
        </nav>
      </aside>
      <div className="ml-60 flex-1">
        <header className="flex items-center justify-between border-b border-gray-200 bg-white px-8 py-4">
          <div className="text-sm text-gray-500">
            欢迎回来，<span className="font-semibold text-gray-800">{user?.nickname || user?.username}</span>
          </div>
          <div className="flex items-center gap-4 text-sm">
            <span className="rounded-full bg-amber-100 px-3 py-1 font-medium text-amber-700">🔥 {user?.streak_days ?? 0} 天</span>
            <span className="rounded-full bg-brand-50 px-3 py-1 font-medium text-brand-700">⭐ {user?.points ?? 0} 分</span>
            <button onClick={handleLogout} className="rounded-lg border border-gray-300 px-3 py-1.5 hover:bg-gray-100">
              退出登录
            </button>
          </div>
        </header>
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
