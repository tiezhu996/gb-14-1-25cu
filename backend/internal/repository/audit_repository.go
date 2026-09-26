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

// 审计仓储哨兵错误。
var ErrAuditNotFound = errors.New("audit log not found")

// AuditRepository 操作审计日志仓储。
type AuditRepository struct {
	coll *mongo.Collection
}

// NewAuditRepository 构造审计日志仓储。
func NewAuditRepository(db *mongo.Database) *AuditRepository {
	return &AuditRepository{coll: db.Collection("audit_logs")}
}

// EnsureIndexes 创建查询索引。
func (r *AuditRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "entity", Value: 1}}},
	})
	return err
}

// Create 写入审计日志。
func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = time.Now()
	if _, err := r.coll.InsertOne(ctx, log); err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

// List 分页查询审计日志（用户/实体/动作过滤）。
func (r *AuditRepository) List(ctx context.Context, filter bson.M, skip, limit int64) ([]*model.AuditLog, int64, error) {
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	defer cur.Close(ctx)
	out := make([]*model.AuditLog, 0)
	for cur.Next(ctx) {
		var l model.AuditLog
		if err := cur.Decode(&l); err != nil {
			return nil, 0, fmt.Errorf("decode audit log: %w", err)
		}
		out = append(out, &l)
	}
	return out, total, nil
}
