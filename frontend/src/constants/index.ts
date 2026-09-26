// 前端常量：与后端 internal/constants 枚举一一对应
export const ROLES = {
  student: 'student',
  admin: 'admin',
} as const

export const ROLE_LABELS: Record<string, string> = {
  student: '学生',
  admin: '管理员',
}

export const DIFFICULTY = {
  easy: 'easy',
  medium: 'medium',
  hard: 'hard',
} as const

export const DIFFICULTY_LABELS: Record<string, string> = {
  easy: '简单',
  medium: '中等',
  hard: '困难',
}

export const DIFFICULTY_CLASSES: Record<string, string> = {
  easy: 'bg-emerald-100 text-emerald-700',
  medium: 'bg-amber-100 text-amber-700',
  hard: 'bg-rose-100 text-rose-700',
}

export const CONTENT_STATUS = {
  draft: 'draft',
  published: 'published',
  archived: 'archived',
} as const

export const CONTENT_STATUS_LABELS: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档',
}

export const SUBMISSION_STATUS = {
  pending: 'pending',
  judging: 'judging',
  accepted: 'accepted',
  partial: 'partial',
  runtime_error: 'runtime_error',
  timeout: 'timeout',
} as const

export const SUBMISSION_STATUS_LABELS: Record<string, string> = {
  pending: '排队中',
  judging: '评测中',
  accepted: '通过',
  partial: '部分通过',
  runtime_error: '运行错误',
  timeout: '超时',
}

export const SUBMISSION_STATUS_CLASSES: Record<string, string> = {
  accepted: 'bg-emerald-100 text-emerald-700',
  partial: 'bg-amber-100 text-amber-700',
  runtime_error: 'bg-rose-100 text-rose-700',
  timeout: 'bg-rose-100 text-rose-700',
  pending: 'bg-gray-100 text-gray-700',
  judging: 'bg-blue-100 text-blue-700',
}

export const LANGUAGES = {
  python: 'python',
  javascript: 'javascript',
  java: 'java',
} as const

export const LANGUAGE_LABELS: Record<string, string> = {
  python: 'Python',
  javascript: 'JavaScript',
  java: 'Java',
}

export const ACHIEVEMENTS = {
  first_ac: 'first_ac',
  streak_7: 'streak_7',
  solved_10: 'solved_10',
  solved_100: 'solved_100',
  first_hard: 'first_hard',
} as const

export const LEADERBOARD_PERIODS = ['daily', 'weekly', 'total'] as const
export const LEADERBOARD_PERIOD_LABELS: Record<string, string> = {
  daily: '日榜',
  weekly: '周榜',
  total: '总榜',
}
