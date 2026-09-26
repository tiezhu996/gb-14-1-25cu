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

// 题目仓储哨兵错误。
var (
	ErrProblemNotFound = errors.New("problem not found")
	ErrProblemExists   = errors.New("problem already exists")
)

// ProblemRepository 题目仓储。
type ProblemRepository struct {
	coll *mongo.Collection
}

// NewProblemRepository 构造题目仓储。
func NewProblemRepository(db *mongo.Database) *ProblemRepository {
	return &ProblemRepository{coll: db.Collection("problems")}
}

// EnsureIndexes 创建标题唯一索引。
func (r *ProblemRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: 1}}, Options: options.Index().SetUnique(true),
	})
	return err
}

// Create 创建题目。
func (r *ProblemRepository) Create(ctx context.Context, p *model.Problem) error {
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, p)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("create problem: %w", ErrProblemExists)
	}
	if err != nil {
		return fmt.Errorf("create problem: %w", err)
	}
	return nil
}

// FindByID 按 ID 查找题目。
func (r *ProblemRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Problem, error) {
	var p model.Problem
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find problem by id: %w", ErrProblemNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find problem by id: %w", err)
	}
	return &p, nil
}

// FindByTitle 按标题查找题目。
func (r *ProblemRepository) FindByTitle(ctx context.Context, title string) (*model.Problem, error) {
	var p model.Problem
	err := r.coll.FindOne(ctx, bson.M{"title": title}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find problem by title: %w", ErrProblemNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find problem by title: %w", err)
	}
	return &p, nil
}

// Update 更新题目内容。
func (r *ProblemRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("update problem: %w", ErrProblemExists)
	}
	if err != nil {
		return fmt.Errorf("update problem: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update problem: %w", ErrProblemNotFound)
	}
	return nil
}

// UpdateStatus 题目状态流转。
func (r *ProblemRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update problem status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update problem status: %w", ErrProblemNotFound)
	}
	return nil
}

// Delete 删除题目。
func (r *ProblemRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete problem: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("delete problem: %w", ErrProblemNotFound)
	}
	return nil
}

// IncSubmit 原子增加提交数。
func (r *ProblemRepository) IncSubmit(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"submit_count": 1}})
	if err != nil {
		return fmt.Errorf("inc problem submit count: %w", err)
	}
	return nil
}

// IncAccepted 原子增加通过数。
func (r *ProblemRepository) IncAccepted(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"accepted_count": 1}})
	if err != nil {
		return fmt.Errorf("inc problem accepted count: %w", err)
	}
	return nil
}

// List 分页查询题目（状态/难度/标签过滤）。
func (r *ProblemRepository) List(ctx context.Context, filter bson.M, skip, limit int64) ([]*model.Problem, int64, error) {
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count problems: %w", err)
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, fmt.Errorf("list problems: %w", err)
	}
	defer cur.Close(ctx)
	problems := make([]*model.Problem, 0)
	for cur.Next(ctx) {
		var p model.Problem
		if err := cur.Decode(&p); err != nil {
			return nil, 0, fmt.Errorf("decode problem: %w", err)
		}
		problems = append(problems, &p)
	}
	return problems, total, nil
}
