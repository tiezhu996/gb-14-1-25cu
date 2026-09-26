// 我的提交记录
import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useSubmissionStore } from '../../stores/submissionStore'
import StatusBadge from '../../components/StatusBadge'
import Pagination from '../../components/Pagination'
import { usePagination } from '../../hooks/usePagination'
import { LANGUAGE_LABELS } from '../../constants'

export default function Submissions() {
  const { submissions, total, loading, fetchSubmissions } = useSubmissionStore()
  const { page, pageSize, setTotal, onPageChange } = usePagination(10)

  useEffect(() => {
    fetchSubmissions({ page, page_size: pageSize }).then(() => {})
  }, [fetchSubmissions, page, pageSize])

  useEffect(() => {
    setTotal(total)
  }, [total, setTotal])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-800">我的提交</h1>
        <p className="mt-1 text-sm text-gray-500">查看每次提交的评测结果与运行详情</p>
      </div>

      {loading && submissions.length === 0 ? (
        <div className="py-20 text-center text-gray-400">加载中...</div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white">
          <table className="min-w-full divide-y divide-gray-200 text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">题目</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">语言</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">状态</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">得分</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">耗时</th>
                <th className="px-5 py-3 text-left font-semibold text-gray-600">提交时间</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {submissions.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50">
                  <td className="px-5 py-3">
                    <Link to={`/problems/${s.problem_id}`} className="font-medium text-brand-600 hover:underline">
                      {s.problem_title}
                    </Link>
                  </td>
                  <td className="px-5 py-3 text-gray-600">{LANGUAGE_LABELS[s.language] || s.language}</td>
                  <td className="px-5 py-3">
                    <StatusBadge value={s.status} kind="submission" />
                  </td>
                  <td className="px-5 py-3 font-medium text-gray-700">{s.score}%</td>
                  <td className="px-5 py-3 text-gray-600">{s.runtime_ms} ms</td>
                  <td className="px-5 py-3 text-gray-500">{s.created_at}</td>
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
