package constants

// 题目难度枚举：同时出现在 model.Problem.Difficulty、DTO、service 积分加权、
// handler 校验、util/formatters.go、前端 constants/difficulty.ts、错误码、日志模板中。
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)

// DifficultyPoints 难度加权积分（排行榜/成就计算依据）。
var DifficultyPoints = map[string]int{
	DifficultyEasy:   10,
	DifficultyMedium: 25,
	DifficultyHard:   50,
}

// ValidDifficulty 校验难度是否合法。
func ValidDifficulty(d string) bool {
	_, ok := DifficultyPoints[d]
	return ok
}
