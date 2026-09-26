// 代码块展示组件（Markdown 代码高亮/讨论帖代码）
import { useState } from 'react'

interface Props {
  code: string
  language?: string
}

export default function CodeBlock({ code, language }: Props) {
  const [copied, setCopied] = useState(false)
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      // 剪贴板不可用时静默
    }
  }
  return (
    <div className="relative rounded-lg bg-gray-900 text-gray-100">
      <div className="flex items-center justify-between border-b border-gray-700 px-4 py-2 text-xs">
        <span className="text-gray-400">{language || 'code'}</span>
        <button onClick={copy} className="rounded bg-gray-700 px-2 py-0.5 hover:bg-gray-600">
          {copied ? '已复制' : '复制'}
        </button>
      </div>
      <pre className="overflow-x-auto p-4 text-sm leading-6">
        <code>{code}</code>
      </pre>
    </div>
  )
}
