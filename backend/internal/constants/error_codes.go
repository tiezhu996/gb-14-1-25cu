// Package constants 集中维护错误码、日志模板、文案与业务枚举。
// 错误码约定：0 表示成功；1xxx 通用；2xxx 用户模块；3xxx 课程模块；
// 4xxx 题目模块；5xxx 提交/评测模块；6xxx 讨论模块；7xxx 成就/排行榜；8xxx 审计/系统。
package constants

// 通用错误码
const (
	CodeOK                = 0
	CodeBadRequest        = 1001
	CodeUnauthorized      = 1002
	CodeForbidden         = 1003
	CodeNotFound          = 1004
	CodeConflict          = 1005
	CodeInternal          = 1006
	CodeValidation        = 1007
	CodeRateLimited       = 1008
	CodeRequestTimeout    = 1009
	CodeInvalidToken      = 1010
	CodeTokenExpired      = 1011
	CodePasswordIncorrect = 1012
)

// 用户模块错误码
const (
	CodeUserNotFound   = 2001
	CodeUserExists     = 2002
	CodeUserBanned     = 2003
	CodeUserSignInLock = 2004
	CodeUserRoleChange = 2005
)

// 课程模块错误码
const (
	CodeCourseNotFound = 3001
	CodeCourseExists   = 3002
	CodeCourseLocked   = 3003
	CodeCourseStatus   = 3004
	CodeCourseOwner    = 3005
)

// 题目模块错误码
const (
	CodeProblemNotFound = 4001
	CodeProblemExists   = 4002
	CodeProblemLocked   = 4003
	CodeProblemNoCases  = 4004
	CodeProblemOwner    = 4005
)

// 提交/评测模块错误码
const (
	CodeSubmissionNotFound = 5001
	CodeSubmissionDenied   = 5002
	CodeJudgeTimeout       = 5003
	CodeJudgeLanguage      = 5004
	CodeJudgeUnavailable   = 5005
)

// 讨论模块错误码
const (
	CodeDiscussionNotFound = 6001
	CodeDiscussionDenied   = 6002
	CodeDiscussionLocked   = 6003
)

// 成就/排行榜模块错误码
const (
	CodeAchievementNotFound = 7001
	CodeLeaderboardEmpty    = 7002
)

// 审计/系统模块错误码
const (
	CodeAuditNotFound = 8001
)
