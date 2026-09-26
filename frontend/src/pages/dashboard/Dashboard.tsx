// 个人学习仪表盘：统计卡片 + 语言分布饼图 + 每日学习热力图
import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useDashboardStore } from '../../stores/dashboardStore'
import { useAuthStore } from '../../stores/authStore'
import StatCard from '../../components/StatCard'
import PieChart from '../../components/PieChart'
import Heatmap from '../../components/Heatmap'
import { formatDuration } from '../../utils/format'
import { signIn } from '../../api/auth'

export default function Dashboard() {
  const { data, loading, fetch } = useDashboardStore()
  const user = useAuthStore((s) => s.user)

  useEffect(() => {
    fetch()
  }, [fetch])

  const doSignIn = async () => {
    try {
      await signIn()
      await fetch()
      useAuthStore.getState().refreshMe()
    } catch {
      // 已签到等错误静默
    }
  }

  if (loading && !data) {
    return <div className="py-20 text-center text-gray-400">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">学习仪表盘</h1>
          <p className="mt-1 text-sm text-gray-500">你好，{user?.nickname || user?.username}，继续加油！</p>
        </div>
        <button onClick={doSignIn} className="rounded-lg bg-amber-500 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-600">
          🔥 每日签到
        </button>
      </div>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-5">
        <StatCard label="累计学习时长" value={formatDuration(data?.total_learning_min || 0)} icon="⏱️" />
        <StatCard label="完成课程" value={data?.completed_courses || 0} icon="📚" />
        <StatCard label="解题总数" value={data?.solved_count || 0} icon="🧩" />
        <StatCard label="提交次数" value={data?.total_submissions || 0} icon="📝" />
        <StatCard label="累计积分" value={data?.points || 0} icon="⭐" />
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 font-semibold text-gray-800">各语言解题分布</h2>
          <PieChart data={data?.language_dist || {}} />
        </div>
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 font-semibold text-gray-800">每日学习热力图（近 90 天）</h2>
          <Heatmap data={data?.daily_activity || {}} />
        </div>
      </div>

      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h2 className="mb-4 font-semibold text-gray-800">快捷入口</h2>
        <div className="grid gap-4 sm:grid-cols-3">
          <Link to="/courses" className="rounded-lg border border-gray-200 p-4 text-center transition hover:border-brand-400 hover:shadow">
            <div className="text-3xl">📚</div>
            <div className="mt-2 font-medium text-gray-800">继续学习课程</div>
            <div className="text-xs text-gray-500">系统化学习路径</div>
          </Link>
          <Link to="/problems" className="rounded-lg border border-gray-200 p-4 text-center transition hover:border-brand-400 hover:shadow">
            <div className="text-3xl">🧩</div>
            <div className="mt-2 font-medium text-gray-800">刷题挑战</div>
            <div className="text-xs text-gray-500">ACM 风格在线评测</div>
          </Link>
          <Link to="/leaderboard" className="rounded-lg border border-gray-200 p-4 text-center transition hover:border-brand-400 hover:shadow">
            <div className="text-3xl">🏆</div>
            <div className="mt-2 font-medium text-gray-800">查看排行榜</div>
            <div className="text-xs text-gray-500">与同学一较高下</div>
          </Link>
        </div>
      </div>
    </div>
  )
}
