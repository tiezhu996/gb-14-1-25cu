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

// 成就仓储哨兵错误。
var (
	ErrAchievementNotFound = errors.New("achievement not found")
	ErrAchievementGranted  = errors.New("achievement already granted")
)

// AchievementRepository 成就徽章仓储。
type AchievementRepository struct {
	coll       *mongo.Collection
	userColl   *mongo.Collection
}

// NewAchievementRepository 构造成就仓储。
func NewAchievementRepository(db *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		coll:     db.Collection("achievements"),
		userColl: db.Collection("user_achievements"),
	}
}

// EnsureIndexes 创建唯一索引（用户+成就码防重复授予）。
func (r *AchievementRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.userColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// Seed 初始化内置成就定义。
func (r *AchievementRepository) Seed(ctx context.Context, defs []struct {
	Code        string
	Name        string
	Description string
	Icon        string
}) error {
	for _, def := range defs {
		var existing model.Achievement
		err := r.coll.FindOne(ctx, bson.M{"code": def.Code}).Decode(&existing)
		if err == nil {
			continue
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("seed achievement: %w", err)
		}
		a := &model.Achievement{
			Code:        def.Code,
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			CreatedAt:   time.Now(),
		}
		a.ID = primitive.NewObjectID()
		if _, err := r.coll.InsertOne(ctx, a); err != nil {
			return fmt.Errorf("seed achievement: %w", err)
		}
	}
	return nil
}

// ListAll 查询全部成就定义。
func (r *AchievementRepository) ListAll(ctx context.Context) ([]*model.Achievement, error) {
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"code": 1}))
	if err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}
	defer cur.Close(ctx)
	out := make([]*model.Achievement, 0)
	for cur.Next(ctx) {
		var a model.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, fmt.Errorf("decode achievement: %w", err)
		}
		out = append(out, &a)
	}
	return out, nil
}

// FindByCode 按成就码查找定义。
func (r *AchievementRepository) FindByCode(ctx context.Context, code string) (*model.Achievement, error) {
	var a model.Achievement
	err := r.coll.FindOne(ctx, bson.M{"code": code}).Decode(&a)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find achievement by code: %w", ErrAchievementNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find achievement by code: %w", err)
	}
	return &a, nil
}

// Grant 授予用户成就（唯一索引保证幂等）。
func (r *AchievementRepository) Grant(ctx context.Context, userID primitive.ObjectID, a *model.Achievement) (*model.UserAchievement, error) {
	ua := &model.UserAchievement{
		UserID:        userID,
		AchievementID: a.ID,
		Code:          a.Code,
		Name:          a.Name,
		Icon:          a.Icon,
		Description:   a.Description,
		EarnedAt:      time.Now(),
	}
	ua.ID = primitive.NewObjectID()
	if _, err := r.userColl.InsertOne(ctx, ua); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("grant achievement: %w", ErrAchievementGranted)
		}
		return nil, fmt.Errorf("grant achievement: %w", err)
	}
	return ua, nil
}

// ListByUser 查询用户已获得的成就。
func (r *AchievementRepository) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]*model.UserAchievement, error) {
	cur, err := r.userColl.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.M{"earned_at": -1}))
	if err != nil {
		return nil, fmt.Errorf("list user achievements: %w", err)
	}
	defer cur.Close(ctx)
	out := make([]*model.UserAchievement, 0)
	for cur.Next(ctx) {
		var ua model.UserAchievement
		if err := cur.Decode(&ua); err != nil {
			return nil, fmt.Errorf("decode user achievement: %w", err)
		}
		out = append(out, &ua)
	}
	return out, nil
}

// HasCode 判断用户是否已获得某成就。
func (r *AchievementRepository) HasCode(ctx context.Context, userID primitive.ObjectID, code string) (bool, error) {
	n, err := r.userColl.CountDocuments(ctx, bson.M{"user_id": userID, "code": code})
	if err != nil {
		return false, fmt.Errorf("check user achievement: %w", err)
	}
	return n > 0, nil
}
