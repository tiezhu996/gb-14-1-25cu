package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo 封装 MongoDB 客户端与数据库句柄。
type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongo 建立 MongoDB 连接并校验连通性。
func NewMongo(ctx context.Context, uri, dbName string, logger *slog.Logger) (*Mongo, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	logger.Info("mongo connected", "db", dbName)
	return &Mongo{Client: client, DB: client.Database(dbName)}, nil
}

// NewRedis 建立 Redis 连接并校验连通性。
func NewRedis(ctx context.Context, addr, password string, logger *slog.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	logger.Info("redis connected", "addr", addr)
	return rdb, nil
}

// Close 关闭 MongoDB 连接。
func (m *Mongo) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
