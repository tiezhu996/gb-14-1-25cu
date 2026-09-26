// 题目管理（管理员）：创建/发布/归档/删除（含测试用例）
import { useEffect, useState } from 'react'
import { useProblemStore } from '../../stores/problemStore'
import { updateProblemStatus, deleteProblem } from '../../api/problem'
import DataTable, { type Column } from '../../components/DataTable'
import StatusBadge from '../../components/StatusBadge'
import ConfirmDialog from '../../components/ConfirmDialog'
import type { Problem, TestCase } from '../../types'

export default function AdminProblems() {
  const { problems, total, loading, fetchProblems } = useProblemStore()
  const [confirm, setConfirm] = useState<Problem | null>(null)
  const [msg, setMsg] = useState('')
  const [form, setForm] = useState({
    title: '',
    description: '',
    difficulty: 'easy',
    languages: 'python',
    tags: '',
    test_cases: [{ input: '1 2', output: '3' }] as TestCase[],
  })

  const load = () => fetchProblems({ page: 1, page_size: 50, status: '' }).then(() => {})
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
      await useProblemStore.getState().createProblem({
        title: form.title,
        description: form.description,
        difficulty: form.difficulty,
        languages: form.languages.split(',').map((s) => s.trim()).filter(Boolean),
        tags: form.tags.split(',').map((s) => s.trim()).filter(Boolean),
        test_cases: form.test_cases.filter((t) => t.input !== '' && t.output !== ''),
        status: 'draft',
      })
      setForm({ ...form, title: '', description: '', tags: '' })
      setMsg('创建成功')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const setStatus = async (p: Problem, status: string) => {
    try {
      await updateProblemStatus(p.id, status)
      setMsg('状态已更新')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const doDelete = async () => {
    if (!confirm) return
    try {
      await deleteProblem(confirm.id)
      setConfirm(null)
      setMsg('已删除')
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const columns: Column<Problem>[] = [
    { key: 'title', title: '标题', render: (p) => <span className="font-medium text-gray-800">{p.title}</span> },
    { key: 'difficulty', title: '难度', render: (p) => <StatusBadge value={p.difficulty} kind="difficulty" /> },
    { key: 'status', title: '状态', render: (p) => <StatusBadge value={p.status} kind="content" /> },
    { key: 'points', title: '积分' },
    { key: 'test_cases', title: '用例数', render: (p) => p.test_cases?.length || 0 },
    {
      key: 'actions',
      title: '操作',
      render: (p) => (
        <div className="flex gap-2">
          {p.status !== 'published' && (
            <button onClick={() => setStatus(p, 'published')} className="rounded bg-emerald-600 px-2.5 py-1 text-xs text-white hover:bg-emerald-700">
              发布
            </button>
          )}
          {p.status === 'published' && (
            <button onClick={() => setStatus(p, 'archived')} className="rounded bg-gray-500 px-2.5 py-1 text-xs text-white hover:bg-gray-600">
              归档
            </button>
          )}
          <button onClick={() => setConfirm(p)} className="rounded bg-rose-600 px-2.5 py-1 text-xs text-white hover:bg-rose-700">
            删除
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">题目管理</h1>
      <div className="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
        <h2 className="font-semibold text-gray-800">新建题目</h2>
        <div className="mt-3 grid gap-3 sm:grid-cols-4">
          <input value={form.title} onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))} placeholder="题目标题" className="rounded-lg border border-gray-300 px-3 py-2 text-sm" />
          <select value={form.difficulty} onChange={(e) => setForm((f) => ({ ...f, difficulty: e.target.value }))} className="rounded-lg border border-gray-300 px-3 py-2 text-sm">
            <option value="easy">简单</option>
            <option value="medium">中等</option>
            <option value="hard">困难</option>
          </select>
          <input value={form.languages} onChange={(e) => setForm((f) => ({ ...f, languages: e.target.value }))} placeholder="语言 python,javascript,java" className="rounded-lg border border-gray-300 px-3 py-2 text-sm" />
          <input value={form.tags} onChange={(e) => setForm((f) => ({ ...f, tags: e.target.value }))} placeholder="标签 逗号分隔" className="rounded-lg border border-gray-300 px-3 py-2 text-sm" />
        </div>
        <textarea value={form.description} onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))} rows={3} placeholder="题目描述（Markdown）" className="mt-3 w-full rounded-lg border border-gray-300 px-3 py-2 text-sm" />
        <div className="mt-3 space-y-2">
          {form.test_cases.map((tc, i) => (
            <div key={i} className="flex gap-2">
              <input
                value={tc.input}
                onChange={(e) => {
                  const next = [...form.test_cases]
                  next[i] = { ...next[i], input: e.target.value }
                  setForm((f) => ({ ...f, test_cases: next }))
                }}
                placeholder={`用例 ${i + 1} 输入`}
                className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm font-mono"
              />
              <input
                value={tc.output}
                onChange={(e) => {
                  const next = [...form.test_cases]
                  next[i] = { ...next[i], output: e.target.value }
                  setForm((f) => ({ ...f, test_cases: next }))
                }}
                placeholder="期望输出"
                className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm font-mono"
              />
              <button
                onClick={() => setForm((f) => ({ ...f, test_cases: f.test_cases.filter((_, idx) => idx !== i) }))}
                className="rounded bg-gray-200 px-3 text-gray-600 hover:bg-gray-300"
              >
                ✕
              </button>
            </div>
          ))}
          <button
            onClick={() => setForm((f) => ({ ...f, test_cases: [...f.test_cases, { input: '', output: '' }] }))}
            className="rounded-lg border border-dashed border-gray-300 px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-50"
          >
            + 添加测试用例
          </button>
        </div>
        <div className="mt-3 flex items-center justify-between">
          {msg && <span className="text-sm text-gray-500">{msg}</span>}
          <button onClick={create} className="ml-auto rounded-lg bg-brand-600 px-4 py-2 text-sm font-semibold text-white hover:bg-brand-700">
            创建题目
          </button>
        </div>
      </div>
      <DataTable columns={columns} rows={problems} rowKey={(p) => p.id} loading={loading} />
      <ConfirmDialog
        open={!!confirm}
        title="删除题目"
        message={`确定要删除题目「${confirm?.title}」吗？该操作不可恢复。`}
        onConfirm={doDelete}
        onCancel={() => setConfirm(null)}
      />
      <p className="text-xs text-gray-400">共 {total} 道题目</p>
    </div>
  )
}
