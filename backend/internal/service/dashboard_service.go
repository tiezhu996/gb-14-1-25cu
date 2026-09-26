package service

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// DashboardService 个人学习仪表盘：统计 + 语言分布 + 热力图。
type DashboardService struct {
	statRepo *repository.UserStatRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewDashboardService 构造仪表盘服务。
func NewDashboardService(statRepo *repository.UserStatRepository, userRepo *repository.UserRepository, logger *slog.Logger) *DashboardService {
	return &DashboardService{statRepo: statRepo, userRepo: userRepo, logger: logger}
}

// Get 获取当前用户学习统计。
func (s *DashboardService) Get(ctx context.Context, userID primitive.ObjectID) (*dto.DashboardResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
	}
	stat, err := s.statRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	resp := &dto.DashboardResponse{
		TotalLearningMin: stat.TotalLearningMin,
		CompletedCourses: stat.CompletedCourses,
		SolvedCount:      user.SolvedCount,
		TotalSubmissions: stat.TotalSubmissions,
		StreakDays:       user.StreakDays,
		Points:           user.Points,
		LanguageDist:     stat.LanguageDist,
		DailyActivity:    s.heatmap(stat.DailyActivity),
	}
	if resp.LanguageDist == nil {
		resp.LanguageDist = map[string]int64{}
	}
	s.logger.Info(constants.LogDashboardQueried, "user_id", userID.Hex())
	return resp, nil
}

// heatmap 生成近 90 天每日活跃热力图（类似 GitHub 贡献图）。
func (s *DashboardService) heatmap(activity map[string]int64) map[string]int64 {
	out := map[string]int64{}
	now := time.Now()
	for i := 89; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		key := util.SignInDailyKey(day)
		if v, ok := activity[key]; ok {
			out[key] = v
		} else {
			out[key] = 0
		}
	}
	return out
}
