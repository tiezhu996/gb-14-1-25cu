package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Achievement 成就徽章定义。
type Achievement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Icon        string             `bson:"icon" json:"icon"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// UserAchievement 用户已获得的成就记录（user_id + code 唯一索引防重复）。
type UserAchievement struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	AchievementID primitive.ObjectID `bson:"achievement_id" json:"achievement_id"`
	Code          string             `bson:"code" json:"code"`
	Name          string             `bson:"name" json:"name"`
	Icon          string             `bson:"icon" json:"icon"`
	Description   string             `bson:"description" json:"description"`
	EarnedAt      time.Time          `bson:"earned_at" json:"earned_at"`
}
