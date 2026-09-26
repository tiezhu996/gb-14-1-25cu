package constants

// 内容状态机枚举：课程与题目共用 draft -> published -> archived。
// 同时出现在 model、DTO、service 状态机、handler 校验、util/formatters.go、
// 前端 constants/status.ts、错误码、日志模板中。
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

// ValidContentStatus 校验内容状态是否合法。
func ValidContentStatus(s string) bool {
	return s == StatusDraft || s == StatusPublished || s == StatusArchived
}

// 用户账号状态枚举。
const (
	UserStatusActive = "active"
	UserStatusBanned = "banned"
)

// ValidUserStatus 校验用户状态。
func ValidUserStatus(s string) bool {
	return s == UserStatusActive || s == UserStatusBanned
}
