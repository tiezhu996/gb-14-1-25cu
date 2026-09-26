// 排行榜：日榜/周榜/总榜
import { useEffect, useState } from 'react'
import { useLeaderboardStore } from '../../stores/leaderboardStore'
import { LEADERBOARD_PERIOD_LABELS, LEADERBOARD_PERIODS } from '../../constants'

export default function Leaderboard() {
  const { entries, loading, fetch } = useLeaderboardStore()
  const [period, setPeriod] = useState('total')

  useEffect(() => {
    fetch(period).then(() => {})
  }, [period, fetch])

  const medal = (rank: number) => (rank === 1 ? '🥇' : rank === 2 ? '🥈' : rank === 3 ? '🥉' : `${rank}`)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">排行榜</h1>
          <p className="mt-1 text-sm text-gray-500">按解题数量与难度加权积分排名</p>
        </div>
        <div className="flex gap-2">
          {LEADERBOARD_PERIODS.map((p) => (
            <button
              key={p}
              onClick={() => setPeriod(p)}
              className={`rounded-lg px-4 py-1.5 text-sm ${period === p ? 'bg-brand-600 text-white' : 'border border-gray-300 hover:bg-gray-100'}`}
            >
              {LEADERBOARD_PERIOD_LABELS[p]}
            </button>
          ))}
        </div>
      </div>

      {loading && entries.length === 0 ? (
        <div className="py-20 text-center text-gray-400">加载中...</div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
          <table className="min-w-full divide-y divide-gray-200 text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">排名</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">用户</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">积分</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">解题数</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {entries.map((e) => (
                <tr key={e.user_id} className="hover:bg-gray-50">
                  <td className="px-5 py-3 text-lg">{medal(e.rank)}</td>
                  <td className="px-5 py-3">
                    <span className="font-medium text-gray-800">{e.nickname}</span>
                    <span className="ml-2 text-xs text-gray-400">@{e.username}</span>
                  </td>
                  <td className="px-5 py-3 font-semibold text-amber-600">{e.points} 分</td>
                  <td className="px-5 py-3 text-gray-600">{e.solved} 题</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
