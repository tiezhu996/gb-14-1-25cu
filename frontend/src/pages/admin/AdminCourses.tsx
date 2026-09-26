// 课程管理（管理员）：创建/发布/归档/删除
import { useEffect, useState } from 'react'
import { useCourseStore } from '../../stores/courseStore'
import { updateCourseStatus, deleteCourse } from '../../api/course'
import DataTable, { type Column } from '../../components/DataTable'
import StatusBadge from '../../components/StatusBadge'
import ConfirmDialog from '../../components/ConfirmDialog'
import type { Course } from '../../types'

export default function AdminCourses() {
  const { courses, total, loading, fetchCourses } = useCourseStore()
  const [form, setForm] = useState({ title: '', description: '', difficulty: 'easy' })
  const [confirm, setConfirm] = useState<Course | null>(null)
  const [msg, setMsg] = useState('')

  const load = () => fetchCourses({ page: 1, page_size: 50, status: '' }).then(() => {})
  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const create = async () => {
    if (!form.title.trim() || !form.description.trim()) {
      setMsg('请填写标题与描述')
      return
    }
    try {
      await useCourseStore.getState().createCourse({
        title: form.title,
        description: form.description,
        difficulty: form.difficulty,
        chapters: [],
        status: 'draft',
      })
      setForm({ title: '', description: '', difficulty: 'easy' })
      setMsg('创建成功')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const setStatus = async (c: Course, status: string) => {
    try {
      await updateCourseStatus(c.id, status)
      setMsg('状态已更新')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const doDelete = async () => {
    if (!confirm) return
    try {
      await deleteCourse(confirm.id)
      setConfirm(null)
      setMsg('已删除')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const columns: Column<Course>[] = [
    { key: 'title', title: '标题', render: (c) => <span className="font-medium text-gray-800">{c.title}</span> },
    { key: 'difficulty', title: '难度', render: (c) => <StatusBadge value={c.difficulty} kind="difficulty" /> },
    { key: 'status', title: '状态', render: (c) => <StatusBadge value={c.status} kind="content" /> },
    { key: 'chapters', title: '章节', render: (c) => c.chapters?.length || 0 },
    {
      key: 'actions',
      title: '操作',
      render: (c) => (
        <div className="flex gap-2">
          {c.status !== 'published' && (
            <button onClick={() => setStatus(c, 'published')} className="rounded bg-emerald-600 px-2.5 py-1 text-xs text-white hover:bg-emerald-700">
              发布
            </button>
          )}
          {c.status === 'published' && (
            <button onClick={() => setStatus(c, 'archived')} className="rounded bg-gray-500 px-2.5 py-1 text-xs text-white hover:bg-gray-600">
              归档
            </button>
          )}
          <button onClick={() => setConfirm(c)} className="rounded bg-rose-600 px-2.5 py-1 text-xs text-white hover:bg-rose-700">
            删除
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">课程管理</h1>
      <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
        <h2 className="font-semibold text-gray-800">新建课程</h2>
        <div className="mt-3 grid gap-3 sm:grid-cols-4">
          <input value={form.title} onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))} placeholder="课程标题" className="rounded-lg border border-gray-300 px-3 py-2 text-sm" />
          <input value={form.description} onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))} placeholder="课程描述" className="rounded-lg border border-gray-300 px-3 py-2 text-sm sm:col-span-2" />
          <select value={form.difficulty} onChange={(e) => setForm((f) => ({ ...f, difficulty: e.target.value }))} className="rounded-lg border border-gray-300 px-3 py-2 text-sm">
            <option value="easy">简单</option>
            <option value="medium">中等</option>
            <option value="hard">困难</option>
          </select>
        </div>
        <div className="mt-3 flex items-center justify-between">
          {msg && <span className="text-sm text-gray-500">{msg}</span>}
          <button onClick={create} className="ml-auto rounded-lg bg-brand-600 px-4 py-2 text-sm font-semibold text-white hover:bg-brand-700">
            创建课程
          </button>
        </div>
      </div>
      <DataTable columns={columns} rows={courses} rowKey={(c) => c.id} loading={loading} />
      <ConfirmDialog
        open={!!confirm}
        title="删除课程"
        message={`确定要删除课程「${confirm?.title}」吗？该操作不可恢复。`}
        onConfirm={doDelete}
        onCancel={() => setConfirm(null)}
      />
      <p className="text-xs text-gray-400">共 {total} 门课程</p>
    </div>
  )
}
