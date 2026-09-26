package constants

// 成就徽章枚举：同时出现在 model.Achievement、service 成就检查、
// 前端 constants/achievements.ts、日志模板、README 枚举清单中。
const (
	AchievementFirstAC    = "first_ac"     // 首次通过题目
	AchievementStreak7    = "streak_7"     // 连续签到 7 天
	AchievementSolved10   = "solved_10"    // 完成 10 题
	AchievementSolved100  = "solved_100"   // 完成 100 题
	AchievementFirstHard  = "first_hard"   // 首次通过困难题
)

// AchievementDefs 内置成就定义（启动时种子化）。
var AchievementDefs = []struct {
	Code        string
	Name        string
	Description string
	Icon        string
}{
	{AchievementFirstAC, "初出茅庐", "首次通过一道编程题目", "🎯"},
	{AchievementStreak7, "持之以恒", "连续签到 7 天", "🔥"},
	{AchievementSolved10, "小试牛刀", "累计通过 10 道题目", "⭐"},
	{AchievementSolved100, "百炼成钢", "累计通过 100 道题目", "🏆"},
	{AchievementFirstHard, "攻坚克难", "首次通过困难难度题目", "🚀"},
}
