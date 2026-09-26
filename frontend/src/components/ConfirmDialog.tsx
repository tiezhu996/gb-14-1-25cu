// 确认对话框（跨页面复用）
interface Props {
  open: boolean
  title: string
  message: string
  confirmText?: string
  onConfirm: () => void
  onCancel: () => void
}

export default function ConfirmDialog({ open, title, message, confirmText = '确认', onConfirm, onCancel }: Props) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" onClick={onCancel}>
      <div className="w-full max-w-md rounded-xl bg-white p-6 shadow-xl" onClick={(e) => e.stopPropagation()}>
        <h3 className="text-lg font-bold text-gray-800">{title}</h3>
        <p className="mt-2 text-sm text-gray-600">{message}</p>
        <div className="mt-6 flex justify-end gap-3">
          <button onClick={onCancel} className="rounded-lg border border-gray-300 px-4 py-2 text-sm hover:bg-gray-100">
            取消
          </button>
          <button onClick={onConfirm} className="rounded-lg bg-rose-600 px-4 py-2 text-sm text-white hover:bg-rose-700">
            {confirmText}
          </button>
        </div>
      </div>
    </div>
  )
}
