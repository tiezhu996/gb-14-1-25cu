// 分页组件（跨页面复用）
interface Props {
  page: number
  pageSize: number
  total: number
  onChange: (page: number) => void
}

export default function Pagination({ page, pageSize, total, onChange }: Props) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize))
  if (totalPages <= 1) return null
  const pages: number[] = []
  for (let p = 1; p <= totalPages; p++) {
    if (p === 1 || p === totalPages || Math.abs(p - page) <= 1) {
      pages.push(p)
    }
  }
  const items: (number | '...')[] = []
  let prev = 0
  for (const p of pages) {
    if (p - prev > 1) items.push('...')
    items.push(p)
    prev = p
  }
  return (
    <div className="mt-4 flex items-center justify-end gap-1 text-sm">
      <button disabled={page <= 1} onClick={() => onChange(page - 1)} className="rounded-md border border-gray-300 px-3 py-1 disabled:opacity-40 hover:bg-gray-100">
        上一页
      </button>
      {items.map((it, i) =>
        it === '...' ? (
          <span key={`e${i}`} className="px-1 text-gray-400">
            ...
          </span>
        ) : (
          <button
            key={it}
            onClick={() => onChange(it)}
            className={`rounded-md px-3 py-1 ${it === page ? 'bg-brand-600 text-white' : 'border border-gray-300 hover:bg-gray-100'}`}
          >
            {it}
          </button>
        ),
      )}
      <button disabled={page >= totalPages} onClick={() => onChange(page + 1)} className="rounded-md border border-gray-300 px-3 py-1 disabled:opacity-40 hover:bg-gray-100">
        下一页
      </button>
    </div>
  )
}
