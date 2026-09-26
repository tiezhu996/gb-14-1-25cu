package constants

// Messages 集中维护接口返回文案、日志文案、错误提示文案（屎山耦合：多处引用）。
const (
	MsgOK              = "ok"
	MsgSuccess         = "success"
	MsgBadRequest      = "请求参数错误"
	MsgUnauthorized    = "未登录或登录已过期"
	MsgForbidden       = "没有权限执行该操作"
	MsgNotFound        = "资源不存在"
	MsgInternalError   = "服务器内部错误"
	MsgValidationError = "参数校验失败"

	MsgUserNotFound     = "用户不存在"
	MsgUserExists       = "用户名或邮箱已被注册"
	MsgUserBanned       = "该账号已被禁用"
	MsgUserWrongPass    = "用户名或密码错误"
	MsgUserSignedToday  = "今日已签到，请明天再来"
	MsgUserSignInOK     = "签到成功，连续签到 %d 天"
	MsgUserRoleChanged  = "角色已更新为 %s"
	MsgUserStatusChanged = "账号状态已更新为 %s"

	MsgCourseNotFound = "课程不存在"
	MsgCourseExists   = "课程标题已存在"
	MsgCourseLocked   = "课程当前状态不允许该操作"
	MsgCourseStatusOK = "课程状态已变更为 %s"
	MsgCourseOwner    = "仅课程作者或管理员可以操作"
	MsgCourseLearned  = "学习时长已记录 %d 分钟"
	MsgCourseDone     = "课程已标记完成"

	MsgProblemNotFound = "题目不存在"
	MsgProblemExists   = "题目标题已存在"
	MsgProblemLocked   = "题目当前状态不允许该操作"
	MsgProblemNoCases  = "题目必须包含至少一个测试用例"
	MsgProblemOwner    = "仅题目创建者或管理员可以操作"
	MsgProblemStatusOK = "题目状态已变更为 %s"

	MsgSubmissionNotFound = "提交记录不存在"
	MsgSubmissionDenied   = "无权查看该提交记录"
	MsgJudgeLanguage      = "不支持的评测语言 %s"
	MsgJudgeTimeout       = "代码运行超时（超过 %d 秒）"
	MsgJudgeUnavailable   = "评测服务暂不可用"
	MsgJudgeCompileError  = "编译错误: %s"

	MsgDiscussionNotFound = "讨论帖不存在"
	MsgDiscussionDenied   = "无权操作该讨论帖"
	MsgDiscussionLocked   = "讨论帖已被隐藏"
	MsgDiscussionBestOK   = "已设为最佳答案"
	MsgDiscussionHiddenOK = "讨论帖已隐藏"

	MsgAchievementNotFound = "成就徽章不存在"
	MsgLeaderboardEmpty    = "排行榜暂无数据"

	MsgInvalidID      = "无效的 %s ID"
	MsgRequestID      = "请求 ID: %s"
	MsgHealthOK       = "codelearn backend is healthy"
)
