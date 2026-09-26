package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserStat 用户学习统计（原子 $inc 维护，避免并发覆盖）。
type UserStat struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID              primitive.ObjectID `bson:"user_id" json:"user_id"`
	TotalLearningMin    int64              `bson:"total_learning_min" json:"total_learning_min"`
	CompletedCourses    int64              `bson:"completed_courses" json:"completed_courses"`
	TotalSubmissions    int64              `bson:"total_submissions" json:"total_submissions"`
	AcceptedSubmissions int64              `bson:"accepted_submissions" json:"accepted_submissions"`
	// LanguageDist 各语言解题分布：language -> count。
	LanguageDist map[string]int64 `bson:"language_dist" json:"language_dist"`
	// DailyActivity 每日学习热力图：yyyy-MM-dd -> 活跃次数。
	DailyActivity map[string]int64 `bson:"daily_activity" json:"daily_activity"`
	CreatedAt     time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at" json:"updated_at"`
}
