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

// 用户仓储哨兵错误。
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserRepository 用户仓储。
type UserRepository struct {
	coll *mongo.Collection
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

// EnsureIndexes 创建唯一索引（并发注册防重）。
func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return err
}

// Create 创建用户；唯一索引冲突时返回 ErrUserExists。
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("create user: %w", ErrUserExists)
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// FindByUsername 按用户名查找。
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find user by username: %w", ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &u, nil
}

// FindByEmail 按邮箱查找。
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find user by email: %w", ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

// FindByID 按 ID 查找。
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var u model.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find user by id: %w", ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

// UpdateProfile 更新昵称/头像/邮箱。
func (r *UserRepository) UpdateProfile(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("update user profile: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update user profile: %w", ErrUserNotFound)
	}
	return nil
}

// UpdateStatus 更新用户状态（active/banned）。
func (r *UserRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update user status: %w", ErrUserNotFound)
	}
	return nil
}

// UpdateRole 更新用户角色（RBAC）。
func (r *UserRepository) UpdateRole(ctx context.Context, id primitive.ObjectID, role string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": role, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update user role: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update user role: %w", ErrUserNotFound)
	}
	return nil
}

// AddPoints 原子增加积分。
func (r *UserRepository) AddPoints(ctx context.Context, id primitive.ObjectID, delta int64) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"points": delta}, "$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("add user points: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("add user points: %w", ErrUserNotFound)
	}
	return nil
}

// MarkSolved 原子增加解题数。
func (r *UserRepository) MarkSolved(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"solved_count": 1}, "$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("mark user solved: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("mark user solved: %w", ErrUserNotFound)
	}
	return nil
}

// UpdateSignIn 更新签到状态：streak 与 last_sign_in_at。
func (r *UserRepository) UpdateSignIn(ctx context.Context, id primitive.ObjectID, streak int, now time.Time) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id},
		bson.M{"$set": bson.M{"streak_days": streak, "last_sign_in_at": now, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update user sign in: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update user sign in: %w", ErrUserNotFound)
	}
	return nil
}

// List 分页查询用户（支持角色/状态过滤）。
func (r *UserRepository) List(ctx context.Context, filter bson.M, skip, limit int64) ([]*model.User, int64, error) {
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer cur.Close(ctx)
	users := make([]*model.User, 0)
	for cur.Next(ctx) {
		var u model.User
		if err := cur.Decode(&u); err != nil {
			return nil, 0, fmt.Errorf("decode user: %w", err)
		}
		users = append(users, &u)
	}
	return users, total, nil
}

// Count 统计用户数量。
func (r *UserRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return n, nil
}
