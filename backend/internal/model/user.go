package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户实体：角色 student/admin，状态 active/banned。
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Nickname     string             `bson:"nickname" json:"nickname"`
	Avatar       string             `bson:"avatar" json:"avatar"`
	Role         string             `bson:"role" json:"role"`
	Status       string             `bson:"status" json:"status"`
	Points       int64              `bson:"points" json:"points"`
	SolvedCount  int64              `bson:"solved_count" json:"solved_count"`
	StreakDays   int                `bson:"streak_days" json:"streak_days"`
	LastSignInAt *time.Time         `bson:"last_sign_in_at" json:"last_sign_in_at"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
