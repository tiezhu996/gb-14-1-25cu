// 通用数据表格组件（跨页面复用）
import type { ReactNode } from 'react'
import EmptyState from './EmptyState'

export interface Column<T> {
  key: string
  title: string
  render?: (row: T) => ReactNode
}

interface Props<T> {
  columns: Column<T>[]
  rows: T[]
  rowKey: (row: T) => string
  loading?: boolean
}

export default function DataTable<T>({ columns, rows, rowKey, loading }: Props<T>) {
  if (loading) {
    return <div className="py-12 text-center text-gray-400">加载中...</div>
  }
  if (rows.length === 0) {
    return <EmptyState />
  }
  return (
    <div className="overflow-x-auto rounded-xl border border-gray-200 bg-white">
      <table className="min-w-full divide-y divide-gray-200 text-sm">
        <thead className="bg-gray-50">
          <tr>
            {columns.map((col) => (
              <th key={col.key} className="px-4 py-3 text-left font-semibold text-gray-600">
                {col.title}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {rows.map((row) => (
            <tr key={rowKey(row)} className="hover:bg-gray-50">
              {columns.map((col) => (
                <td key={col.key} className="px-4 py-3">
                  {col.render ? col.render(row) : String((row as Record<string, unknown>)[col.key] ?? '-')}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
