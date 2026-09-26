// 空状态占位组件（跨页面复用）
interface Props {
  title?: string
  description?: string
  icon?: string
}

export default function EmptyState({ title = '暂无数据', description, icon = '📭' }: Props) {
  return (
    <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-gray-300 bg-white py-16 text-center">
      <div className="text-5xl">{icon}</div>
      <h3 className="mt-4 text-lg font-semibold text-gray-700">{title}</h3>
      {description && <p className="mt-1 max-w-md text-sm text-gray-500">{description}</p>}
    </div>
  )
}
