package repository

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/blueship581/codelearn/internal/model"
)

func testMongo(t *testing.T) *mongo.Database {
	t.Helper()
	uri := os.Getenv("TEST_MONGO_URI")
	if uri == "" {
		t.Skip("TEST_MONGO_URI not set, skip integration test")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("mongo connect: %v", err)
	}
	t.Cleanup(func() { _ = client.Disconnect(context.Background()) })
	return client.Database("codelearn_repo_test")
}

func TestUserRepositoryCRUD(t *testing.T) {
	db := testMongo(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("ensure indexes: %v", err)
	}
	_, _ = db.Collection("users").DeleteMany(ctx, map[string]any{"username": "repo_test_user"})

	u := &model.User{Username: "repo_test_user", Email: "repo_test_user@example.com", PasswordHash: "hash", Role: "student", Status: "active"}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if u.ID.IsZero() {
		t.Error("expected generated id")
	}

	found, err := repo.FindByUsername(ctx, "repo_test_user")
	if err != nil {
		t.Fatalf("FindByUsername error: %v", err)
	}
	if found.Email != u.Email {
		t.Errorf("email = %q", found.Email)
	}

	// 唯一索引冲突
	dup := &model.User{Username: "repo_test_user", Email: "other@example.com", PasswordHash: "hash"}
	if err := repo.Create(ctx, dup); err == nil {
		t.Error("expected duplicate key error")
	}

	if err := repo.UpdateStatus(ctx, u.ID, "banned"); err != nil {
		t.Fatalf("UpdateStatus error: %v", err)
	}
	if err := repo.AddPoints(ctx, u.ID, 10); err != nil {
		t.Fatalf("AddPoints error: %v", err)
	}
	if err := repo.MarkSolved(ctx, u.ID); err != nil {
		t.Fatalf("MarkSolved error: %v", err)
	}
	after, err := repo.FindByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if after.Points != 10 || after.SolvedCount != 1 || after.Status != "banned" {
		t.Errorf("after = %+v", after)
	}

	_, total, err := repo.List(ctx, map[string]any{"username": "repo_test_user"}, 0, 10)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}

	_, _ = db.Collection("users").DeleteMany(ctx, map[string]any{"username": "repo_test_user"})
}
