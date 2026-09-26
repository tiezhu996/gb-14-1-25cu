// 课程详情：Markdown 内容 + 章节列表 + 学习记录/完成
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { getCourse, recordLearn, completeCourse } from '../../api/course'
import type { Course } from '../../types'
import StatusBadge from '../../components/StatusBadge'

export default function CourseDetail() {
  const { id } = useParams()
  const [course, setCourse] = useState<Course | null>(null)
  const [activeChapter, setActiveChapter] = useState(0)
  const [msg, setMsg] = useState('')

  useEffect(() => {
    if (id) {
      getCourse(id).then(setCourse).catch(() => setCourse(null))
      // 记录 1 分钟学习时长（演示）
      recordLearn(id, 1).catch(() => undefined)
    }
  }, [id])

  const complete = async () => {
    if (!id) return
    try {
      await completeCourse(id)
      setMsg('🎉 课程已标记完成！')
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  if (!course) {
    return <div className="py-20 text-center text-gray-400">加载中...</div>
  }

  const chapter = course.chapters?.[activeChapter]

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      <div className="lg:col-span-2">
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="flex items-center gap-2">
            <StatusBadge value={course.difficulty} kind="difficulty" />
            <StatusBadge value={course.status} kind="content" />
          </div>
          <h1 className="mt-3 text-2xl font-bold text-gray-800">{course.title}</h1>
          <p className="mt-2 text-sm text-gray-500">{course.description}</p>
          <div className="mt-6 markdown-body">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{chapter ? chapter.content : course.markdown || ''}</ReactMarkdown>
          </div>
        </div>
      </div>
      <div>
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h2 className="font-semibold text-gray-800">课程章节</h2>
          {course.chapters?.length ? (
            <ul className="mt-4 space-y-2">
              {course.chapters.map((ch, i) => (
                <li key={ch.title}>
                  <button
                    onClick={() => setActiveChapter(i)}
                    className={`w-full rounded-lg border px-4 py-3 text-left text-sm transition ${
                      activeChapter === i ? 'border-brand-500 bg-brand-50' : 'border-gray-200 hover:bg-gray-50'
                    }`}
                  >
                    <div className="font-medium text-gray-800">
                      {i + 1}. {ch.title}
                    </div>
                    <div className="mt-1 text-xs text-gray-400">预计 {ch.duration} 分钟</div>
                  </button>
                </li>
              ))}
            </ul>
          ) : (
            <p className="mt-4 text-sm text-gray-400">暂无章节</p>
          )}
          <button onClick={complete} className="mt-6 w-full rounded-lg bg-emerald-600 py-2 text-sm font-semibold text-white hover:bg-emerald-700">
            ✅ 标记完成
          </button>
          {msg && <div className="mt-3 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{msg}</div>}
        </div>
      </div>
    </div>
  )
}
