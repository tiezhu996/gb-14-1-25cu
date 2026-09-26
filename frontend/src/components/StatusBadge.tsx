// 通用状态徽标：按枚举值映射颜色与文案
import { DIFFICULTY_CLASSES, DIFFICULTY_LABELS, CONTENT_STATUS_LABELS, SUBMISSION_STATUS_CLASSES, SUBMISSION_STATUS_LABELS, ROLE_LABELS } from '../constants'

interface Props {
  value: string
  kind?: 'difficulty' | 'content' | 'submission' | 'role' | 'default'
}

export default function StatusBadge({ value, kind = 'default' }: Props) {
  let label = value
  let className = 'bg-gray-100 text-gray-700'
  if (kind === 'difficulty') {
    label = DIFFICULTY_LABELS[value] || value
    className = DIFFICULTY_CLASSES[value] || className
  } else if (kind === 'content') {
    label = CONTENT_STATUS_LABELS[value] || value
    className = value === 'published' ? 'bg-emerald-100 text-emerald-700' : value === 'archived' ? 'bg-gray-200 text-gray-600' : 'bg-amber-100 text-amber-700'
  } else if (kind === 'submission') {
    label = SUBMISSION_STATUS_LABELS[value] || value
    className = SUBMISSION_STATUS_CLASSES[value] || className
  } else if (kind === 'role') {
    label = ROLE_LABELS[value] || value
    className = value === 'admin' ? 'bg-violet-100 text-violet-700' : 'bg-sky-100 text-sky-700'
  }
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  )
}
