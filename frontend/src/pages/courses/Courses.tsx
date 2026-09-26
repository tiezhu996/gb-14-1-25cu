// 课程中心：课程列表（按难度筛选）
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useCourseStore } from '../../stores/courseStore'
import StatusBadge from '../../components/StatusBadge'
import Pagination from '../../components/Pagination'
import { usePagination } from '../../hooks/usePagination'
import { DIFFICULTY } from '../../constants'

export default function Courses() {
  const { courses, total, loading, fetchCourses } = useCourseStore()
  const { page, pageSize, setTotal, onPageChange } = usePagination(9)
  const [difficulty, setDifficulty] = useState('')

  useEffect(() => {
    fetchCourses({ page, page_size: pageSize, difficulty: difficulty || undefined }).then(() => {})
  }, [fetchCourses, page, pageSize, difficulty])

  useEffect(() => {
    setTotal(total)
  }, [total, setTotal])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">课程中心</h1>
          <p className="mt-1 text-sm text-gray-500">系统化的学习路径，从入门到进阶</p>
        </div>
        <div className="flex gap-2">
          {['', DIFFICULTY.easy, DIFFICULTY.medium, DIFFICULTY.hard].map((d) => (
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

      {loading && courses.length === 0 ? (
        <div className="py-20 text-center text-gray-400">加载中...</div>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {courses.map((course) => (
            <Link
              key={course.id}
              to={`/courses/${course.id}`}
              className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
            >
              <div className="h-36 bg-cover bg-center" style={{ backgroundImage: `url(${course.cover})` }} />
              <div className="p-5">
                <div className="flex items-center justify-between">
                  <StatusBadge value={course.difficulty} kind="difficulty" />
                  <StatusBadge value={course.status} kind="content" />
                </div>
                <h2 className="mt-3 text-lg font-semibold text-gray-800">{course.title}</h2>
                <p className="mt-1 line-clamp-2 text-sm text-gray-500">{course.description}</p>
                <div className="mt-3 flex items-center justify-between text-xs text-gray-400">
                  <span>{course.chapters?.length || 0} 个章节</span>
                  <span>{course.author_name}</span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      <Pagination page={page} pageSize={pageSize} total={total} onChange={onPageChange} />
    </div>
  )
}
