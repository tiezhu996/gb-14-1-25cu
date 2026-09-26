package service

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// testMongo 连接测试库；未配置 TEST_MONGO_URI 时跳过。
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
	return client.Database("codelearn_test")
}

func TestUserServiceRegisterLogin(t *testing.T) {
	db := testMongo(t)
	logger := util.NewLogger("error")
	userRepo := repository.NewUserRepository(db)
	statRepo := repository.NewUserStatRepository(db)
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		t.Fatalf("ensure indexes: %v", err)
	}
	svc := NewUserService(userRepo, statRepo, logger, "test-secret")
	ctx := context.Background()

	// 清理可能存在的历史测试用户
	_, _ = db.Collection("users").DeleteMany(ctx, map[string]any{"username": "svc_test_user"})

	register := &dto.RegisterRequest{Username: "svc_test_user", Email: "svc_test_user@example.com", Password: "secret123", Nickname: "测试用户"}
	resp, err := svc.Register(ctx, register)
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token after register")
	}

	// 重复注册应报用户已存在
	if _, err := svc.Register(ctx, register); err == nil {
		t.Error("expected error on duplicate register")
	}

	// 登录成功
	login, err := svc.Login(ctx, &dto.LoginRequest{Username: "svc_test_user", Password: "secret123"})
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if login.Token == "" {
		t.Error("expected token after login")
	}
}
