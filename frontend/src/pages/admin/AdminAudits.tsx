// 审计日志（管理员）：查看全站写操作记录
import { useEffect, useState } from 'react'
import { useAuditStore } from '../../stores/auditStore'
import DataTable, { type Column } from '../../components/DataTable'
import Pagination from '../../components/Pagination'
import { usePagination } from '../../hooks/usePagination'
import type { AuditLog } from '../../types'

export default function AdminAudits() {
  const { logs, total, loading, fetch } = useAuditStore()
  const { page, pageSize, setTotal, onPageChange } = usePagination(20)
  const [entity, setEntity] = useState('')

  useEffect(() => {
    fetch({ page, page_size: pageSize, entity: entity || undefined }).then(() => {})
  }, [fetch, page, pageSize, entity])

  useEffect(() => {
    setTotal(total)
  }, [total, setTotal])

  const columns: Column<AuditLog>[] = [
    { key: 'created_at', title: '时间', render: (l) => <span className="text-gray-500">{l.created_at}</span> },
    { key: 'username', title: '用户', render: (l) => <span className="font-medium text-gray-800">{l.username || '匿名'}</span> },
    { key: 'method', title: '方法', render: (l) => <code className="rounded bg-gray-100 px-1.5 py-0.5 text-xs">{l.method}</code> },
    { key: 'path', title: '路径', render: (l) => <code className="text-xs text-gray-600">{l.path}</code> },
    { key: 'entity', title: '实体' },
    { key: 'ip', title: 'IP', render: (l) => <span className="text-gray-500">{l.ip}</span> },
    { key: 'request_id', title: 'Request ID', render: (l) => <code className="text-xs text-gray-500">{l.request_id.slice(0, 8)}</code> },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">审计日志</h1>
          <p className="mt-1 text-sm text-gray-500">全站写操作记录（由审计中间件自动采集）</p>
        </div>
        <select value={entity} onChange={(e) => { setEntity(e.target.value); onPageChange(1) }} className="rounded-lg border border-gray-300 px-3 py-1.5 text-sm">
          <option value="">全部实体</option>
          <option value="auth">auth</option>
          <option value="users">users</option>
          <option value="courses">courses</option>
          <option value="problems">problems</option>
          <option value="submissions">submissions</option>
          <option value="discussions">discussions</option>
        </select>
      </div>
      <DataTable columns={columns} rows={logs} rowKey={(l) => l.id} loading={loading} />
      <Pagination page={page} pageSize={pageSize} total={total} onChange={onPageChange} />
    </div>
  )
}
