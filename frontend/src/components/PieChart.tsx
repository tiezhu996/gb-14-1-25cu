// SVG 环形饼图组件（语言分布）
import { LANGUAGE_LABELS } from '../constants'

interface Props {
  data: Record<string, number>
}

const COLORS = ['#3380ff', '#10b981', '#f59e0b', '#8b5cf6', '#ef4444']

export default function PieChart({ data }: Props) {
  const entries = Object.entries(data).filter(([, v]) => v > 0)
  const total = entries.reduce((sum, [, v]) => sum + v, 0)
  if (total === 0) {
    return <div className="flex h-48 items-center justify-center text-sm text-gray-400">暂无解题数据</div>
  }
  let angle = 0
  const arcs = entries.map(([lang, value], i) => {
    const start = angle
    const sweep = (value / total) * 360
    angle += sweep
    const large = sweep > 180 ? 1 : 0
    const startRad = ((start - 90) * Math.PI) / 180
    const endRad = ((start + sweep - 90) * Math.PI) / 180
    const x1 = 50 + 38 * Math.cos(startRad)
    const y1 = 50 + 38 * Math.sin(startRad)
    const x2 = 50 + 38 * Math.cos(endRad)
    const y2 = 50 + 38 * Math.sin(endRad)
    return {
      key: lang,
      path: `M 50 50 L ${x1} ${y1} A 38 38 0 ${large} 1 ${x2} ${y2} Z`,
      color: COLORS[i % COLORS.length],
      percent: Math.round((value / total) * 100),
    }
  })
  return (
    <div className="flex items-center gap-6">
      <svg viewBox="0 0 100 100" className="h-40 w-40">
        {arcs.map((a) => (
          <path key={a.key} d={a.path} fill={a.color} stroke="#fff" strokeWidth="1" />
        ))}
        <circle cx="50" cy="50" r="24" fill="#fff" />
      </svg>
      <div className="space-y-2">
        {entries.map(([lang, value], i) => (
          <div key={lang} className="flex items-center gap-2 text-sm">
            <span className="h-3 w-3 rounded-sm" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
            <span className="text-gray-600">{LANGUAGE_LABELS[lang] || lang}</span>
            <span className="font-medium text-gray-800">{value} 题</span>
          </div>
        ))}
      </div>
    </div>
  )
}
