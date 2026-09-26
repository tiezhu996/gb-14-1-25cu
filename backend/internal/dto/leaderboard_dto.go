package dto

// LeaderboardEntryResponse 排行榜条目。
type LeaderboardEntryResponse struct {
	Rank     int    `json:"rank"`
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Points   int64  `json:"points"`
	Solved   int64  `json:"solved"`
}

// DashboardResponse 个人学习仪表盘统计。
type DashboardResponse struct {
	TotalLearningMin   int64            `json:"total_learning_min"`
	CompletedCourses   int64            `json:"completed_courses"`
	SolvedCount        int64            `json:"solved_count"`
	TotalSubmissions   int64            `json:"total_submissions"`
	StreakDays         int              `json:"streak_days"`
	Points             int64            `json:"points"`
	LanguageDist       map[string]int64 `json:"language_dist"`
	DailyActivity      map[string]int64 `json:"daily_activity"`
}
