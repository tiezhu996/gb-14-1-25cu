package constants

// 提交评测状态机：pending -> judging -> accepted | partial | runtime_error | timeout。
// 同时出现在 model.Submission.Status、DTO、service 状态机、judge service、
// util/formatters.go、前端 constants/submission.ts、错误码、日志模板中。
const (
	SubmissionPending     = "pending"
	SubmissionJudging     = "judging"
	SubmissionAccepted    = "accepted"
	SubmissionPartial     = "partial"
	SubmissionRuntimeError = "runtime_error"
	SubmissionTimeout     = "timeout"
)

// DefaultJudgeTimeout 评测默认超时（秒），与提示词"运行超时限制 10 秒"对齐。
const DefaultJudgeTimeout = 10

// 支持的评测语言枚举。
const (
	LanguagePython    = "python"
	LanguageJavaScript = "javascript"
	LanguageJava      = "java"
)

// SupportedLanguages 返回全部支持的评测语言。
func SupportedLanguages() []string {
	return []string{LanguagePython, LanguageJavaScript, LanguageJava}
}

// ValidLanguage 校验评测语言是否支持。
func ValidLanguage(lang string) bool {
	for _, l := range SupportedLanguages() {
		if l == lang {
			return true
		}
	}
	return false
}

// ValidSubmissionStatus 校验提交状态。
func ValidSubmissionStatus(s string) bool {
	switch s {
	case SubmissionPending, SubmissionJudging, SubmissionAccepted,
		SubmissionPartial, SubmissionRuntimeError, SubmissionTimeout:
		return true
	}
	return false
}
