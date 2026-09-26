package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/blueship581/codelearn/internal/constants"
)

// Formatters 集中提供日期、状态文本、类型文本等格式化逻辑（屎山耦合：多处引用）。

// FormatDateTime 格式化时间为 yyyy-MM-dd HH:mm:ss。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 格式化时间为 yyyy-MM-dd。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatRoleText 角色文本（数据库角色字段 -> 前端徽标）。
func FormatRoleText(role string) string {
	switch role {
	case constants.RoleAdmin:
		return "管理员"
	case constants.RoleStudent:
		return "学生"
	default:
		return "未知"
	}
}

// FormatDifficultyText 难度文本（题目难度 -> 前端徽标）。
func FormatDifficultyText(d string) string {
	switch d {
	case constants.DifficultyEasy:
		return "简单"
	case constants.DifficultyMedium:
		return "中等"
	case constants.DifficultyHard:
		return "困难"
	default:
		return "未知"
	}
}

// FormatDifficultyClass 难度对应的 Tailwind 徽标类名。
func FormatDifficultyClass(d string) string {
	switch d {
	case constants.DifficultyEasy:
		return "bg-emerald-100 text-emerald-700"
	case constants.DifficultyMedium:
		return "bg-amber-100 text-amber-700"
	case constants.DifficultyHard:
		return "bg-rose-100 text-rose-700"
	default:
		return "bg-gray-100 text-gray-700"
	}
}

// FormatContentStatusText 内容状态文本（课程/题目状态机 -> 前端徽标）。
func FormatContentStatusText(s string) string {
	switch s {
	case constants.StatusDraft:
		return "草稿"
	case constants.StatusPublished:
		return "已发布"
	case constants.StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}

// FormatSubmissionStatusText 提交状态文本（评测状态机 -> 前端徽标）。
func FormatSubmissionStatusText(s string) string {
	switch s {
	case constants.SubmissionPending:
		return "排队中"
	case constants.SubmissionJudging:
		return "评测中"
	case constants.SubmissionAccepted:
		return "通过"
	case constants.SubmissionPartial:
		return "部分通过"
	case constants.SubmissionRuntimeError:
		return "运行错误"
	case constants.SubmissionTimeout:
		return "超时"
	default:
		return "未知"
	}
}

// FormatSubmissionStatusClass 提交状态对应的徽标类名。
func FormatSubmissionStatusClass(s string) string {
	switch s {
	case constants.SubmissionAccepted:
		return "bg-emerald-100 text-emerald-700"
	case constants.SubmissionPartial:
		return "bg-amber-100 text-amber-700"
	case constants.SubmissionRuntimeError, constants.SubmissionTimeout:
		return "bg-rose-100 text-rose-700"
	default:
		return "bg-gray-100 text-gray-700"
	}
}

// FormatLanguageText 语言文本（评测语言 -> 前端徽标）。
func FormatLanguageText(lang string) string {
	switch lang {
	case constants.LanguagePython:
		return "Python"
	case constants.LanguageJavaScript:
		return "JavaScript"
	case constants.LanguageJava:
		return "Java"
	default:
		return lang
	}
}

// FormatUserStatusText 用户状态文本。
func FormatUserStatusText(s string) string {
	switch s {
	case constants.UserStatusActive:
		return "正常"
	case constants.UserStatusBanned:
		return "已禁用"
	default:
		return "未知"
	}
}

// FormatDiscussionSortText 讨论排序方式文本。
func FormatDiscussionSortText(s string) string {
	switch s {
	case constants.DiscussionSortBest:
		return "最佳答案"
	case constants.DiscussionSortNew:
		return "最新发布"
	case constants.DiscussionSortVotes:
		return "最多点赞"
	default:
		return "未知"
	}
}

// NormalizeOutput 规范化评测输出（去首尾空白、统一换行、去除每行尾部空格），用于 AC 比对。
func NormalizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// SignInDailyKey 生成签到/热力图的日期键。
func SignInDailyKey(t time.Time) string {
	return t.Format("2006-01-02")
}

// PointsLabel 积分展示文本。
func PointsLabel(p int64) string {
	return fmt.Sprintf("%d 分", p)
}
