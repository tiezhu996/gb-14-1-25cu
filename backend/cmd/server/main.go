package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blueship581/codelearn/internal/config"
	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/database"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/router"
	"github.com/blueship581/codelearn/internal/seed"
	"github.com/blueship581/codelearn/internal/util"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config failed", "error", err.Error())
		os.Exit(1)
	}
	logger := util.NewLogger(cfg.LogLevel)
	logger.Info(constants.LogServerStarted, "env", cfg.Env, "port", cfg.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongo, err := database.NewMongo(ctx, cfg.MongoURI, cfg.DBName, logger)
	if err != nil {
		logger.Error(constants.LogDBConnectFailed, "error", err.Error())
		os.Exit(1)
	}
	defer mongo.Close(ctx)

	rdb, err := database.NewRedis(ctx, cfg.RedisAddr, cfg.RedisPass, logger)
	if err != nil {
		logger.Error(constants.LogRedisConnectFailed, "error", err.Error())
		os.Exit(1)
	}
	defer rdb.Close()

	// 初始化索引与种子数据
	if err := ensureIndexes(ctx, mongo); err != nil {
		logger.Error("ensure indexes failed", "error", err.Error())
		os.Exit(1)
	}
	if cfg.SeedAdmin == "true" {
		if err := seed.Seed(ctx, mongo.DB, logger); err != nil {
			logger.Error("seed data failed", "error", err.Error())
			os.Exit(1)
		}
	}

	engine := router.New(router.Deps{
		DB:           mongo,
		Redis:        rdb,
		Logger:       logger,
		JWTSecret:    cfg.JWTSecret,
		JudgeTimeout: cfg.JudgeTimeOut,
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: engine,
	}

	go func() {
		logger.Info("http server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server error", "error", err.Error())
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(constants.LogServerShutdown)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err.Error())
	}
}

// ensureIndexes 创建全部集合唯一/查询索引。
func ensureIndexes(ctx context.Context, mongo *database.Mongo) error {
	db := mongo.DB
	if err := repository.NewUserRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewCourseRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewProblemRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewSubmissionRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewDiscussionRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewAchievementRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewAuditRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := repository.NewUserStatRepository(db).EnsureIndexes(ctx); err != nil {
		return err
	}
	return nil
}
