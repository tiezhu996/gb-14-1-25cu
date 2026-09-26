package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/blueship581/codelearn/internal/model"
)

// 提交仓储哨兵错误。
var ErrSubmissionNotFound = errors.New("submission not found")

// SubmissionRepository 提交评测仓储。
type SubmissionRepository struct {
	coll *mongo.Collection
}

// NewSubmissionRepository 构造提交仓储。
func NewSubmissionRepository(db *mongo.Database) *SubmissionRepository {
	return &SubmissionRepository{coll: db.Collection("submissions")}
}

// EnsureIndexes 创建常用查询索引（user_id、problem_id）。
func (r *SubmissionRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "problem_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}

// Create 创建提交记录（初始 pending）。
func (r *SubmissionRepository) Create(ctx context.Context, s *model.Submission) error {
	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("create submission: %w", err)
	}
	return nil
}

// FindByID 按 ID 查找提交。
func (r *SubmissionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Submission, error) {
	var s model.Submission
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find submission by id: %w", ErrSubmissionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find submission by id: %w", err)
	}
	return &s, nil
}

// UpdateResult 更新评测结果。
func (r *SubmissionRepository) UpdateResult(ctx context.Context, id primitive.ObjectID, status string, score int, pointsAwarded, runtimeMs int64, results []model.JudgeResult, errMsg string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":         status,
		"score":          score,
		"points_awarded": pointsAwarded,
		"runtime_ms":     runtimeMs,
		"results":        results,
		"error_message":  errMsg,
	}})
	if err != nil {
		return fmt.Errorf("update submission result: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update submission result: %w", ErrSubmissionNotFound)
	}
	return nil
}

// List 分页查询提交（用户/题目/状态过滤）。
func (r *SubmissionRepository) List(ctx context.Context, filter bson.M, skip, limit int64) ([]*model.Submission, int64, error) {
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count submissions: %w", err)
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, fmt.Errorf("list submissions: %w", err)
	}
	defer cur.Close(ctx)
	subs := make([]*model.Submission, 0)
	for cur.Next(ctx) {
		var s model.Submission
		if err := cur.Decode(&s); err != nil {
			return nil, 0, fmt.Errorf("decode submission: %w", err)
		}
		subs = append(subs, &s)
	}
	return subs, total, nil
}

// AggregatePoints 聚合得分：按日期范围、状态统计每个用户的积分与解题数（排行榜复用）。
func (r *SubmissionRepository) AggregatePoints(ctx context.Context, filter bson.M) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id":     "$user_id",
			"points":  bson.M{"$sum": "$points_awarded"},
			"solved":  bson.M{"$sum": 1},
			"nickname": bson.M{"$last": "$username"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "points", Value: -1}, {Key: "solved", Value: -1}}}},
	}
	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate submission points: %w", err)
	}
	defer cur.Close(ctx)
	out := make([]bson.M, 0)
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode aggregate: %w", err)
		}
		out = append(out, doc)
	}
	return out, nil
}

// CountAcceptedByUser 统计用户在某题目的通过次数（成就判定复用）。
func (r *SubmissionRepository) CountAcceptedByUser(ctx context.Context, userID, problemID primitive.ObjectID) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{
		"user_id":    userID,
		"problem_id": problemID,
		"status":     "accepted",
	})
	if err != nil {
		return 0, fmt.Errorf("count accepted by user: %w", err)
	}
	return n, nil
}
