package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chapter 课程章节：标题 + Markdown 讲解 + 预计学习时长（分钟）。
type Chapter struct {
	Title    string `bson:"title" json:"title"`
	Content  string `bson:"content" json:"content"`
	Duration int    `bson:"duration" json:"duration"`
}

// Course 课程实体：状态机 draft -> published -> archived。
type Course struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Cover       string             `bson:"cover" json:"cover"`
	Markdown    string             `bson:"markdown" json:"markdown"`
	Chapters    []Chapter          `bson:"chapters" json:"chapters"`
	Difficulty  string             `bson:"difficulty" json:"difficulty"`
	Status      string             `bson:"status" json:"status"`
	AuthorID    primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorName  string             `bson:"author_name" json:"author_name"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
