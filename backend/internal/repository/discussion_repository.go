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

// 讨论仓储哨兵错误。
var ErrDiscussionNotFound = errors.New("discussion not found")

// DiscussionRepository 讨论帖仓储。
type DiscussionRepository struct {
	coll *mongo.Collection
}

// NewDiscussionRepository 构造讨论帖仓储。
func NewDiscussionRepository(db *mongo.Database) *DiscussionRepository {
	return &DiscussionRepository{coll: db.Collection("discussions")}
}

// EnsureIndexes 创建查询索引。
func (r *DiscussionRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "problem_id", Value: 1}, {Key: "is_best", Value: -1}, {Key: "vote_count", Value: -1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
	})
	return err
}

// Create 发布讨论帖。
func (r *DiscussionRepository) Create(ctx context.Context, d *model.Discussion) error {
	d.ID = primitive.NewObjectID()
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, d)
	if err != nil {
		return fmt.Errorf("create discussion: %w", err)
	}
	return nil
}

// FindByID 按 ID 查找讨论帖。
func (r *DiscussionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Discussion, error) {
	var d model.Discussion
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find discussion by id: %w", ErrDiscussionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find discussion by id: %w", err)
	}
	return &d, nil
}

// ListByProblem 按题目查询讨论帖（排序方式：best/new/votes）。
func (r *DiscussionRepository) ListByProblem(ctx context.Context, problemID primitive.ObjectID, sort string) ([]*model.Discussion, error) {
	sortDoc := bson.D{{Key: "created_at", Value: -1}}
	switch sort {
	case "best":
		sortDoc = bson.D{{Key: "is_best", Value: -1}, {Key: "vote_count", Value: -1}, {Key: "created_at", Value: -1}}
	case "votes":
		sortDoc = bson.D{{Key: "vote_count", Value: -1}, {Key: "created_at", Value: -1}}
	}
	cur, err := r.coll.Find(ctx, bson.M{"problem_id": problemID, "status": "active"},
		options.Find().SetSort(sortDoc))
	if err != nil {
		return nil, fmt.Errorf("list discussions by problem: %w", err)
	}
	defer cur.Close(ctx)
	discs := make([]*model.Discussion, 0)
	for cur.Next(ctx) {
		var d model.Discussion
		if err := cur.Decode(&d); err != nil {
			return nil, fmt.Errorf("decode discussion: %w", err)
		}
		discs = append(discs, &d)
	}
	return discs, nil
}

// Vote 原子投票（幂等：voter_ids 数组防重复）。
func (r *DiscussionRepository) Vote(ctx context.Context, id, userID primitive.ObjectID, delta int) error {
	update := bson.M{"$inc": bson.M{"vote_count": delta}, "$addToSet": bson.M{"voter_ids": userID}, "$set": bson.M{"updated_at": time.Now()}}
	if delta < 0 {
		update = bson.M{"$inc": bson.M{"vote_count": delta}, "$pull": bson.M{"voter_ids": userID}, "$set": bson.M{"updated_at": time.Now()}}
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("vote discussion: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("vote discussion: %w", ErrDiscussionNotFound)
	}
	return nil
}

// HasVoted 判断用户是否已点赞。
func (r *DiscussionRepository) HasVoted(ctx context.Context, id, userID primitive.ObjectID) (bool, error) {
	var d struct {
		VoterIDs []primitive.ObjectID `bson:"voter_ids"`
	}
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"voter_ids": 1})).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, fmt.Errorf("has voted: %w", ErrDiscussionNotFound)
	}
	if err != nil {
		return false, fmt.Errorf("has voted: %w", err)
	}
	for _, vid := range d.VoterIDs {
		if vid == userID {
			return true, nil
		}
	}
	return false, nil
}

// MarkBest 标记最佳答案（先清除同题旧最佳，再设置新最佳）。
func (r *DiscussionRepository) MarkBest(ctx context.Context, problemID, id primitive.ObjectID) error {
	if _, err := r.coll.UpdateMany(ctx, bson.M{"problem_id": problemID, "is_best": true},
		bson.M{"$set": bson.M{"is_best": false, "updated_at": time.Now()}}); err != nil {
		return fmt.Errorf("clear best discussion: %w", err)
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_best": true, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("mark best discussion: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("mark best discussion: %w", ErrDiscussionNotFound)
	}
	return nil
}

// UpdateStatus 隐藏/恢复讨论帖。
func (r *DiscussionRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update discussion status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update discussion status: %w", ErrDiscussionNotFound)
	}
	return nil
}
