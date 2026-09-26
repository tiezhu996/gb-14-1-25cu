// 题目讨论区组件：发布/排序/点赞/最佳答案
import { useEffect, useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { useDiscussionStore } from '../../stores/discussionStore'
import { useAuthStore } from '../../stores/authStore'
import CodeBlock from '../../components/CodeBlock'
import LanguageSelect from '../../components/LanguageSelect'
import { LANGUAGE_LABELS } from '../../constants'

interface Props {
  problemId: string
}

export default function Discussions({ problemId }: Props) {
  const { discussions, loading, fetchDiscussions, createDiscussion, vote, markBest } = useDiscussionStore()
  const user = useAuthStore((s) => s.user)
  const [sort, setSort] = useState('best')
  const [form, setForm] = useState({ title: '', content: '', code: '', language: 'python' })
  const [msg, setMsg] = useState('')

  useEffect(() => {
    fetchDiscussions(problemId, sort).then(() => {})
  }, [problemId, sort, fetchDiscussions])

  const publish = async () => {
    if (!form.title.trim() || !form.content.trim()) {
      setMsg('请填写标题与内容')
      return
    }
    try {
      await createDiscussion(problemId, {
        title: form.title,
        content: form.content,
        code: form.code || undefined,
        language: form.code ? form.language : undefined,
      })
      setForm({ title: '', content: '', code: '', language: 'python' })
      setMsg('发布成功 🎉')
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  return (
    <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex items-center justify-between">
        <h2 className="font-semibold text-gray-800">💬 讨论区</h2>
        <select value={sort} onChange={(e) => setSort(e.target.value)} className="rounded-lg border border-gray-300 px-3 py-1.5 text-sm">
          <option value="best">最佳答案优先</option>
          <option value="new">最新发布</option>
          <option value="votes">最多点赞</option>
        </select>
      </div>

      <div className="mt-4 rounded-lg border border-gray-200 bg-gray-50 p-4">
        <div className="grid gap-3 sm:grid-cols-2">
          <input
            value={form.title}
            onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
            placeholder="标题：如「分享我的双指针解法」"
            className="rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
          <div className="flex items-center gap-2">
            <input
              value={form.code}
              onChange={(e) => setForm((f) => ({ ...f, code: e.target.value }))}
              placeholder="附上代码（选填）"
              className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm"
            />
            <LanguageSelect value={form.language} onChange={(v) => setForm((f) => ({ ...f, language: v }))} />
          </div>
        </div>
        <textarea
          value={form.content}
          onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
          placeholder="支持 Markdown 格式的题解分享与思路讨论..."
          rows={3}
          className="mt-3 w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
        />
        <div className="mt-3 flex items-center justify-between">
          {msg && <span className="text-sm text-gray-500">{msg}</span>}
          <button onClick={publish} className="ml-auto rounded-lg bg-brand-600 px-4 py-2 text-sm font-semibold text-white hover:bg-brand-700">
            发布讨论
          </button>
        </div>
      </div>

      {loading ? (
        <div className="py-8 text-center text-gray-400">加载中...</div>
      ) : discussions.length === 0 ? (
        <div className="py-8 text-center text-gray-400">还没有讨论，来抢沙发吧～</div>
      ) : (
        <ul className="mt-4 space-y-4">
          {discussions.map((d) => (
            <li key={d.id} className="rounded-lg border border-gray-200 p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-gray-800">{d.nickname || d.username}</span>
                  <span className="text-xs text-gray-400">{d.created_at}</span>
                  {d.is_best && (
                    <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700">🏅 最佳答案</span>
                  )}
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <button
                    onClick={() => vote(d.id).catch(() => undefined)}
                    className="rounded-lg border border-gray-300 px-2.5 py-1 hover:bg-gray-100"
                  >
                    👍 {d.vote_count}
                  </button>
                  {user?.role === 'admin' && (
                    <button
                      onClick={() => markBest(problemId, d.id).catch(() => undefined)}
                      className="rounded-lg border border-amber-300 px-2.5 py-1 text-amber-600 hover:bg-amber-50"
                    >
                      🏅 设为最佳
                    </button>
                  )}
                </div>
              </div>
              <h3 className="mt-2 font-semibold text-gray-800">{d.title}</h3>
              <div className="mt-1 markdown-body text-sm">
                <ReactMarkdown remarkPlugins={[remarkGfm]}>{d.content}</ReactMarkdown>
              </div>
              {d.code && <div className="mt-2"><CodeBlock code={d.code} language={LANGUAGE_LABELS[d.language] || d.language} /></div>}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
