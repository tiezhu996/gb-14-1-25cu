package service

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// LeaderboardService 排行榜：日榜/周榜/总榜（聚合提交积分）。
type LeaderboardService struct {
	subRepo *repository.SubmissionRepository
	userRepo *repository.UserRepository
	logger  *slog.Logger
}

// NewLeaderboardService 构造排行榜服务。
func NewLeaderboardService(subRepo *repository.SubmissionRepository, userRepo *repository.UserRepository, logger *slog.Logger) *LeaderboardService {
	return &LeaderboardService{subRepo: subRepo, userRepo: userRepo, logger: logger}
}

// Get 查询排行榜：period 支持 daily/weekly/total。
func (s *LeaderboardService) Get(ctx context.Context, period string, limit int64) ([]dto.LeaderboardEntryResponse, error) {
	filter := bson.M{"status": constants.SubmissionAccepted}
	now := time.Now()
	switch period {
	case "daily":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		filter["created_at"] = bson.M{"$gte": start}
	case "weekly":
		start := now.AddDate(0, 0, -7)
		filter["created_at"] = bson.M{"$gte": start}
	case "total":
		// 总榜不过滤时间
	default:
		return nil, util.WrapAppError(constants.CodeBadRequest, "排行榜周期必须为 daily/weekly/total", nil)
	}
	rows, err := s.subRepo.AggregatePoints(ctx, filter)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if limit <= 0 {
		limit = 50
	}
	out := make([]dto.LeaderboardEntryResponse, 0, len(rows))
	for i, row := range rows {
		if int64(i) >= limit {
			break
		}
		userIDHex, _ := row["_id"].(primitive.ObjectID)
		points, _ := row["points"].(int64)
		solved, _ := row["solved"].(int64)
		username, _ := row["nickname"].(string)
		nickname := username
		if user, err := s.userRepo.FindByID(ctx, userIDHex); err == nil {
			nickname = user.Nickname
			username = user.Username
		}
		out = append(out, dto.LeaderboardEntryResponse{
			Rank:     i + 1,
			UserID:   userIDHex.Hex(),
			Nickname: nickname,
			Username: username,
			Points:   points,
			Solved:   solved,
		})
	}
	s.logger.Info(constants.LogLeaderboardQueried, "period", period, "count", len(out))
	return out, nil
}
