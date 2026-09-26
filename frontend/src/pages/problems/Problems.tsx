// 题库：题目列表（难度/标签筛选）
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useProblemStore } from '../../stores/problemStore'
import StatusBadge from '../../components/StatusBadge'
import Pagination from '../../components/Pagination'
import { usePagination } from '../../hooks/usePagination'

export default function Problems() {
  const { problems, total, loading, fetchProblems } = useProblemStore()
  const { page, pageSize, setTotal, onPageChange } = usePagination(10)
  const [difficulty, setDifficulty] = useState('')

  useEffect(() => {
    fetchProblems({ page, page_size: pageSize, difficulty: difficulty || undefined }).then(() => {})
  }, [fetchProblems, page, pageSize, difficulty])

  useEffect(() => {
    setTotal(total)
  }, [total, setTotal])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">题库练习</h1>
          <p className="mt-1 text-sm text-gray-500">提交代码，在线评测，与全站同学同台竞技</p>
        </div>
        <div className="flex gap-2">
          {['', 'easy', 'medium', 'hard'].map((d) => (
            <button
              key={d || 'all'}
              onClick={() => {
                setDifficulty(d)
                onPageChange(1)
              }}
              className={`rounded-lg px-3 py-1.5 text-sm ${difficulty === d ? 'bg-brand-600 text-white' : 'border border-gray-300 hover:bg-gray-100'}`}
            >
              {d === '' ? '全部' : d === 'easy' ? '简单' : d === 'medium' ? '中等' : '困难'}
            </button>
          ))}
        </div>
      </div>

      {loading && problems.length === 0 ? (
        <div className="py-20 text-center text-gray-400">加载中...</div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
          <table className="min-w-full divide-y divide-gray-200 text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">题目</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">难度</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">积分</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">标签</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">通过率</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {problems.map((p) => (
                <tr key={p.id} className="hover:bg-gray-50">
                  <td className="px-5 py-3">
                    <Link to={`/problems/${p.id}`} className="font-medium text-brand-600 hover:underline">
                      {p.title}
                    </Link>
                  </td>
                  <td className="px-5 py-3">
                    <StatusBadge value={p.difficulty} kind="difficulty" />
                  </td>
                  <td className="px-5 py-3 font-medium text-gray-700">{p.points}</td>
                  <td className="px-5 py-3">
                    <div className="flex flex-wrap gap-1">
                      {p.tags?.map((t) => (
                        <span key={t} className="rounded bg-gray-100 px-2 py-0.5 text-xs text-gray-600">
                          {t}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td className="px-5 py-3 text-gray-600">
                    {p.submit_count > 0 ? `${Math.round((p.accepted_count / p.submit_count) * 100)}%` : '-'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Pagination page={page} pageSize={pageSize} total={total} onChange={onPageChange} />
    </div>
  )
}
