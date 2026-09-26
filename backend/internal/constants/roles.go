package constants

// 角色枚举：同时出现在 model.User.Role、DTO、service 状态机、handler 校验、
// middleware/rbac.go、util/formatters.go、前端 constants/roles.ts、错误码、日志模板中。
const (
	RoleStudent = "student" // 学生（默认）
	RoleAdmin   = "admin"   // 管理员
)

// ValidRole 校验角色是否合法（handler 入参校验使用）。
func ValidRole(role string) bool {
	return role == RoleStudent || role == RoleAdmin
}
