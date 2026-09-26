package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Discussion 题目讨论帖：支持 Markdown + 代码块，可标记最佳答案。
type Discussion struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProblemID  primitive.ObjectID `bson:"problem_id" json:"problem_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username   string             `bson:"username" json:"username"`
	Nickname   string             `bson:"nickname" json:"nickname"`
	Title      string             `bson:"title" json:"title"`
	Content    string             `bson:"content" json:"content"`
	Code       string             `bson:"code" json:"code"`
	Language   string             `bson:"language" json:"language"`
	IsBest     bool               `bson:"is_best" json:"is_best"`
	VoteCount  int                `bson:"vote_count" json:"vote_count"`
	VoterIDs   []primitive.ObjectID `bson:"voter_ids" json:"-"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
