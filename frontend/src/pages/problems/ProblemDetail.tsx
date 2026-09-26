// 题目详情：Monaco 在线 IDE + 提交评测结果 + 讨论区
import { useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import Editor from '@monaco-editor/react'
import * as monaco from 'monaco-editor'
import { loader } from '@monaco-editor/react'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import { getProblem } from '../../api/problem'
import { submitCode } from '../../api/submission'
import type { Problem, Submission } from '../../types'
import StatusBadge from '../../components/StatusBadge'
import LanguageSelect from '../../components/LanguageSelect'
import Discussions from './Discussions'
import { SUBMISSION_STATUS_LABELS } from '../../constants'

// Monaco Worker 本地配置（不依赖 CDN）
self.MonacoEnvironment = {
  getWorker: () => new EditorWorker(),
}
loader.config({ monaco })

const STARTER_CODE: Record<string, string> = {
  python: '# 读取输入并按题意输出结果\na, b = map(int, input().split())\nprint(a + b)',
  javascript: "// 读取输入并按题意输出结果\nconst [a, b] = require('fs').readFileSync(0, 'utf-8').trim().split(/\\s+/).map(Number);\nconsole.log(a + b);",
  java: "import java.util.*;\n\npublic class Main {\n    public static void main(String[] args) {\n        Scanner sc = new Scanner(System.in);\n        int a = sc.nextInt();\n        int b = sc.nextInt();\n        System.out.println(a + b);\n    }\n}",
}

export default function ProblemDetail() {
  const { id } = useParams()
  const [problem, setProblem] = useState<Problem | null>(null)
  const [language, setLanguage] = useState('python')
  const [code, setCode] = useState(STARTER_CODE.python)
  const [submitting, setSubmitting] = useState(false)
  const [submission, setSubmission] = useState<Submission | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    if (id) {
      getProblem(id)
        .then((p) => {
          setProblem(p)
          const lang = p.languages?.[0] || 'python'
          setLanguage(lang)
          setCode(STARTER_CODE[lang] || '')
        })
        .catch((e) => setError((e as Error).message))
    }
  }, [id])

  const monacoOptions = useMemo(
    () => ({
      minimap: { enabled: false },
      fontSize: 14,
      automaticLayout: true,
      scrollBeyondLastLine: false,
      tabSize: 4,
    }),
    [],
  )

  const switchLanguage = (lang: string) => {
    setLanguage(lang)
    setCode(STARTER_CODE[lang] || '')
  }

  const doSubmit = async () => {
    if (!id) return
    setSubmitting(true)
    setSubmission(null)
    setError('')
    try {
      const resp = await submitCode(id, { language, code })
      setSubmission(resp)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setSubmitting(false)
    }
  }

  if (!problem) {
    return <div className="py-20 text-center text-gray-400">{error || '加载中...'}</div>
  }

  return (
    <div className="space-y-6">
      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="mr-2 text-2xl font-bold text-gray-800">{problem.title}</h1>
          <StatusBadge value={problem.difficulty} kind="difficulty" />
          <span className="rounded-full bg-amber-100 px-2.5 py-0.5 text-xs font-medium text-amber-700">⭐ {problem.points} 分</span>
        </div>
        <div className="mt-2 flex flex-wrap gap-1 text-xs text-gray-400">
          {problem.tags?.map((t) => (
            <span key={t} className="rounded bg-gray-100 px-2 py-0.5 text-gray-600">
              {t}
            </span>
          ))}
          <span>时间限制：{problem.time_limit}s</span>
          <span>提交：{problem.submit_count}</span>
          <span>通过：{problem.accepted_count}</span>
        </div>
        <div className="mt-4 markdown-body">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{problem.description}</ReactMarkdown>
        </div>
      </div>

      <div className="rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="flex items-center justify-between border-b border-gray-100 px-6 py-3">
          <div className="flex items-center gap-3">
            <span className="font-semibold text-gray-800">在线 IDE</span>
            <LanguageSelect value={language} onChange={switchLanguage} disabled={submitting} />
          </div>
          <button
            onClick={doSubmit}
            disabled={submitting}
            className="rounded-lg bg-brand-600 px-5 py-2 text-sm font-semibold text-white hover:bg-brand-700 disabled:opacity-60"
          >
            {submitting ? '评测中...' : '🚀 提交评测'}
          </button>
        </div>
        <Editor
          height="420px"
          language={language === 'javascript' ? 'javascript' : language}
          value={code}
          onChange={(v) => setCode(v || '')}
          options={monacoOptions}
          theme="vs"
        />
      </div>

      {submission && (
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <div className="flex items-center gap-3">
            <h2 className="font-semibold text-gray-800">评测结果</h2>
            <StatusBadge value={submission.status} kind="submission" />
            <span className="text-sm text-gray-500">
              {SUBMISSION_STATUS_LABELS[submission.status]} · 得分 {submission.score}% · 耗时 {submission.runtime_ms}ms
            </span>
            {submission.points_awarded > 0 && (
              <span className="rounded-full bg-amber-100 px-2.5 py-0.5 text-xs font-medium text-amber-700">+{submission.points_awarded} 分</span>
            )}
          </div>
          {submission.error_message && (
            <div className="mt-3 rounded-lg bg-rose-50 px-4 py-3 text-sm text-rose-700">{submission.error_message}</div>
          )}
          {submission.results?.length > 0 && (
            <div className="mt-4 overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 text-sm">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-2 text-left font-semibold text-gray-600">用例</th>
                    <th className="px-4 py-2 text-left font-semibold text-gray-600">输入</th>
                    <th className="px-4 py-2 text-left font-semibold text-gray-600">期望输出</th>
                    <th className="px-4 py-2 text-left font-semibold text-gray-600">实际输出</th>
                    <th className="px-4 py-2 text-left font-semibold text-gray-600">结果</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {submission.results.map((r) => (
                    <tr key={r.test_case_index}>
                      <td className="px-4 py-2">#{r.test_case_index + 1}</td>
                      <td className="px-4 py-2 font-mono text-xs">{r.input || '(空)'}</td>
                      <td className="px-4 py-2 font-mono text-xs">{r.expected || '(空)'}</td>
                      <td className="px-4 py-2 font-mono text-xs">{r.actual || '(空)'}</td>
                      <td className="px-4 py-2">
                        {r.passed ? <span className="text-emerald-600">✅ 通过</span> : <span className="text-rose-600">❌ 未通过</span>}
                        {r.error_message && <span className="ml-1 text-xs text-rose-500">{r.error_message}</span>}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {id && <Discussions problemId={id} />}
    </div>
  )
}
