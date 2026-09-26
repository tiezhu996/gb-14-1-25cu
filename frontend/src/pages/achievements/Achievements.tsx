// 成就徽章：全部徽章 + 我的徽章
import { useEffect } from 'react'
import { useAchievementStore } from '../../stores/achievementStore'

export default function Achievements() {
  const { all, mine, fetchAll, fetchMine } = useAchievementStore()
  const mineCodes = new Set(mine.map((m) => m.code))

  useEffect(() => {
    fetchAll().then(() => {})
    fetchMine().then(() => {})
  }, [fetchAll, fetchMine])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-800">成就徽章</h1>
        <p className="mt-1 text-sm text-gray-500">已解锁 {mine.length} / {all.length} 个徽章</p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {all.map((a) => {
          const unlocked = mineCodes.has(a.code)
          return (
            <div
              key={a.id}
              className={`rounded-xl border p-6 text-center shadow-sm transition ${
                unlocked ? 'border-amber-300 bg-amber-50' : 'border-gray-200 bg-white opacity-70'
              }`}
            >
              <div className={`mx-auto flex h-16 w-16 items-center justify-center rounded-full text-4xl ${unlocked ? 'bg-amber-100' : 'bg-gray-100 grayscale'}`}>
                {a.icon}
              </div>
              <h2 className="mt-3 font-semibold text-gray-800">{a.name}</h2>
              <p className="mt-1 text-sm text-gray-500">{a.description}</p>
              {unlocked ? (
                <span className="mt-3 inline-block rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-700">已解锁</span>
              ) : (
                <span className="mt-3 inline-block rounded-full bg-gray-100 px-3 py-1 text-xs text-gray-500">未解锁</span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
