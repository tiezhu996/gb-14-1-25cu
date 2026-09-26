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

// 课程仓储哨兵错误。
var (
	ErrCourseNotFound = errors.New("course not found")
	ErrCourseExists   = errors.New("course already exists")
)

// CourseRepository 课程仓储。
type CourseRepository struct {
	coll *mongo.Collection
}

// NewCourseRepository 构造课程仓储。
func NewCourseRepository(db *mongo.Database) *CourseRepository {
	return &CourseRepository{coll: db.Collection("courses")}
}

// EnsureIndexes 创建标题唯一索引。
func (r *CourseRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: 1}}, Options: options.Index().SetUnique(true),
	})
	return err
}

// Create 创建课程。
func (r *CourseRepository) Create(ctx context.Context, c *model.Course) error {
	c.ID = primitive.NewObjectID()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, c)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("create course: %w", ErrCourseExists)
	}
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	return nil
}

// FindByID 按 ID 查找课程。
func (r *CourseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Course, error) {
	var c model.Course
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find course by id: %w", ErrCourseNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find course by id: %w", err)
	}
	return &c, nil
}

// FindByTitle 按标题查找课程。
func (r *CourseRepository) FindByTitle(ctx context.Context, title string) (*model.Course, error) {
	var c model.Course
	err := r.coll.FindOne(ctx, bson.M{"title": title}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find course by title: %w", ErrCourseNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("find course by title: %w", err)
	}
	return &c, nil
}

// Update 更新课程内容（仅标题唯一冲突时报错）。
func (r *CourseRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("update course: %w", ErrCourseExists)
	}
	if err != nil {
		return fmt.Errorf("update course: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update course: %w", ErrCourseNotFound)
	}
	return nil
}

// UpdateStatus 课程状态流转。
func (r *CourseRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("update course status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("update course status: %w", ErrCourseNotFound)
	}
	return nil
}

// Delete 删除课程。
func (r *CourseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete course: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("delete course: %w", ErrCourseNotFound)
	}
	return nil
}

// List 分页查询课程（状态/难度过滤）。
func (r *CourseRepository) List(ctx context.Context, filter bson.M, skip, limit int64) ([]*model.Course, int64, error) {
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count courses: %w", err)
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, fmt.Errorf("list courses: %w", err)
	}
	defer cur.Close(ctx)
	courses := make([]*model.Course, 0)
	for cur.Next(ctx) {
		var c model.Course
		if err := cur.Decode(&c); err != nil {
			return nil, 0, fmt.Errorf("decode course: %w", err)
		}
		courses = append(courses, &c)
	}
	return courses, total, nil
}
