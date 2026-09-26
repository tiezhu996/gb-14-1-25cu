// 学习热力图组件（类似 GitHub 贡献图，近 90 天）
import { formatHeatLevel } from '../utils/format'

interface Props {
  data: Record<string, number>
}

export default function Heatmap({ data }: Props) {
  const days = Object.keys(data).sort()
  const maxPerRow = 15
  const rows: string[][] = []
  for (let i = 0; i < days.length; i += maxPerRow) {
    rows.push(days.slice(i, i + maxPerRow))
  }
  return (
    <div className="overflow-x-auto">
      <div className="flex gap-1">
        {rows.map((row, ri) => (
          <div key={ri} className="flex flex-col gap-1">
            {row.map((day) => {
              const count = data[day] || 0
              return (
                <div
                  key={day}
                  title={`${day}: ${count} 次活跃`}
                  className={`h-3 w-3 rounded-sm ${formatHeatLevel(count)}`}
                />
              )
            })}
          </div>
        ))}
      </div>
      <div className="mt-2 flex items-center gap-1 text-xs text-gray-400">
        <span>少</span>
        {[0, 1, 2, 3, 4].map((n) => (
          <span key={n} className={`h-3 w-3 rounded-sm ${formatHeatLevel(n)}`} />
        ))}
        <span>多</span>
      </div>
    </div>
  )
}
