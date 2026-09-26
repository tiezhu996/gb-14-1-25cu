// 统计卡片（仪表盘复用）
interface Props {
  label: string
  value: string | number
  icon: string
}

export default function StatCard({ label, value, icon }: Props) {
  return (
    <div className="flex items-center gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
      <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-brand-50 text-2xl">{icon}</div>
      <div>
        <div className="text-2xl font-bold text-gray-800">{value}</div>
        <div className="text-sm text-gray-500">{label}</div>
      </div>
    </div>
  )
}
