// 用户管理（管理员）：禁用/启用、角色调整
import { useEffect, useState } from 'react'
import { listUsers, updateUserStatus, updateUserRole } from '../../api/user'
import DataTable, { type Column } from '../../components/DataTable'
import StatusBadge from '../../components/StatusBadge'
import type { User } from '../../types'

export default function AdminUsers() {
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [msg, setMsg] = useState('')

  const load = async () => {
    setLoading(true)
    try {
      const data = await listUsers({ page: 1, page_size: 50 })
      setUsers(data.list)
      setTotal(data.total)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load().then(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const toggleStatus = async (u: User) => {
    try {
      await updateUserStatus(u.id, u.status === 'active' ? 'banned' : 'active')
      setMsg(`已将 ${u.username} 状态更新`)
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const toggleRole = async (u: User) => {
    try {
      await updateUserRole(u.id, u.role === 'admin' ? 'student' : 'admin')
      setMsg(`已将 ${u.username} 角色更新`)
      load()
    } catch (e) {
      setMsg((e as Error).message)
    }
  }

  const columns: Column<User>[] = [
    { key: 'username', title: '用户名', render: (u) => <span className="font-medium text-gray-800">{u.username}</span> },
    { key: 'nickname', title: '昵称' },
    { key: 'email', title: '邮箱' },
    { key: 'role', title: '角色', render: (u) => <StatusBadge value={u.role} kind="role" /> },
    { key: 'status', title: '状态', render: (u) => <StatusBadge value={u.status} kind="default" /> },
    { key: 'points', title: '积分' },
    { key: 'solved_count', title: '解题数' },
    {
      key: 'actions',
      title: '操作',
      render: (u) => (
        <div className="flex gap-2">
          <button onClick={() => toggleStatus(u)} className="rounded bg-gray-600 px-2.5 py-1 text-xs text-white hover:bg-gray-700">
            {u.status === 'active' ? '禁用' : '启用'}
          </button>
          <button onClick={() => toggleRole(u)} className="rounded bg-violet-600 px-2.5 py-1 text-xs text-white hover:bg-violet-700">
            {u.role === 'admin' ? '降为学生' : '设为管理员'}
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">用户管理</h1>
        {msg && <span className="text-sm text-gray-500">{msg}</span>}
      </div>
      <DataTable columns={columns} rows={users} rowKey={(u) => u.id} loading={loading} />
      <p className="text-xs text-gray-400">共 {total} 位用户</p>
    </div>
  )
}
