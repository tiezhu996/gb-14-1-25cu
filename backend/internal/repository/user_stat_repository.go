package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/blueship581/codelearn/internal/model"
)

// UserStatRepository 用户学习统计仓储（原子 $inc，无事务依赖）。
type UserStatRepository struct {
	coll *mongo.Collection
}

// NewUserStatRepository 构造统计仓储。
func NewUserStatRepository(db *mongo.Database) *UserStatRepository {
	return &UserStatRepository{coll: db.Collection("user_stats")}
}

// EnsureIndexes 创建 user_id 唯一索引。
func (r *UserStatRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true),
	})
	return err
}

// Upsert 初始化用户统计文档（幂等）。
func (r *UserStatRepository) Upsert(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"user_id": userID},
		bson.M{"$setOnInsert": bson.M{
			"user_id":            userID,
			"total_learning_min": 0,
			"completed_courses":  0,
			"total_submissions":  0,
			"accepted_submissions": 0,
			"language_dist":      bson.M{},
			"daily_activity":     bson.M{},
			"created_at":         now,
			"updated_at":         now,
		}},
		options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("upsert user stat: %w", err)
	}
	return nil
}

// AddLearning 原子增加学习时长并记录日活跃。
func (r *UserStatRepository) AddLearning(ctx context.Context, userID primitive.ObjectID, minutes int64, dayKey string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{
		"$inc": bson.M{"total_learning_min": minutes, "daily_activity." + dayKey: 1},
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		return fmt.Errorf("add learning stat: %w", err)
	}
	return nil
}

// CompleteCourse 原子增加完成课程数并记录日活跃。
func (r *UserStatRepository) CompleteCourse(ctx context.Context, userID primitive.ObjectID, dayKey string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{
		"$inc": bson.M{"completed_courses": 1, "daily_activity." + dayKey: 1},
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		return fmt.Errorf("complete course stat: %w", err)
	}
	return nil
}

// AddSubmission 原子增加提交数；accepted=true 时同时增加通过数与语言分布。
func (r *UserStatRepository) AddSubmission(ctx context.Context, userID primitive.ObjectID, language string, accepted bool, dayKey string) error {
	inc := bson.M{"total_submissions": 1, "daily_activity." + dayKey: 1}
	if accepted {
		inc["accepted_submissions"] = 1
		inc["language_dist."+language] = 1
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{
		"$inc": inc,
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		return fmt.Errorf("add submission stat: %w", err)
	}
	return nil
}

// GetByUser 查询用户统计。
func (r *UserStatRepository) GetByUser(ctx context.Context, userID primitive.ObjectID) (*model.UserStat, error) {
	var s model.UserStat
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("get user stat: %w", err)
	}
	return &s, nil
}
