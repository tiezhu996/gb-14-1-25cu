// 通用格式化工具：日期、时长、状态文本
export function formatDateTime(s?: string): string {
  if (!s) return '-'
  return s
}

export function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes} 分钟`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m > 0 ? `${h} 小时 ${m} 分钟` : `${h} 小时`
}

export function formatRuntime(ms: number): string {
  if (ms <= 0) return '-'
  return `${ms} ms`
}

export function formatHeatLevel(count: number): string {
  if (count <= 0) return 'bg-gray-100'
  if (count === 1) return 'bg-emerald-200'
  if (count <= 3) return 'bg-emerald-300'
  if (count <= 6) return 'bg-emerald-400'
  return 'bg-emerald-500'
}
