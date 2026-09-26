package constants

// 讨论帖状态枚举。
const (
	DiscussionActive = "active"
	DiscussionHidden = "hidden"
)

// 讨论帖排序方式枚举。
const (
	DiscussionSortBest  = "best"
	DiscussionSortNew   = "new"
	DiscussionSortVotes = "votes"
)

// ValidDiscussionSort 校验排序方式。
func ValidDiscussionSort(s string) bool {
	return s == DiscussionSortBest || s == DiscussionSortNew || s == DiscussionSortVotes
}

// ValidDiscussionStatus 校验讨论帖状态。
func ValidDiscussionStatus(s string) bool {
	return s == DiscussionActive || s == DiscussionHidden
}
