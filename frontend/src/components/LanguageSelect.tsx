// 评测语言选择（跨页面复用）
import { LANGUAGE_LABELS, LANGUAGES } from '../constants'

interface Props {
  value: string
  onChange: (v: string) => void
  disabled?: boolean
}

export default function LanguageSelect({ value, onChange, disabled }: Props) {
  return (
    <select
      value={value}
      disabled={disabled}
      onChange={(e) => onChange(e.target.value)}
      className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm focus:border-brand-500 focus:outline-none"
    >
      {Object.values(LANGUAGES).map((lang) => (
        <option key={lang} value={lang}>
          {LANGUAGE_LABELS[lang]}
        </option>
      ))}
    </select>
  )
}
