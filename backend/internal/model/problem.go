package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCase ACM 风格测试用例：标准输入 + 期望输出。
type TestCase struct {
	Input  string `bson:"input" json:"input"`
	Output string `bson:"output" json:"output"`
}

// Problem 编程题目实体：难度 easy/medium/hard，状态 draft -> published -> archived。
type Problem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Difficulty    string             `bson:"difficulty" json:"difficulty"`
	Languages     []string           `bson:"languages" json:"languages"`
	Tags          []string           `bson:"tags" json:"tags"`
	TestCases     []TestCase         `bson:"test_cases" json:"test_cases"`
	TimeLimit     int                `bson:"time_limit" json:"time_limit"`
	Status        string             `bson:"status" json:"status"`
	Points        int                `bson:"points" json:"points"`
	AcceptedCount int64              `bson:"accepted_count" json:"accepted_count"`
	SubmitCount   int64              `bson:"submit_count" json:"submit_count"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedByName string             `bson:"created_by_name" json:"created_by_name"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
