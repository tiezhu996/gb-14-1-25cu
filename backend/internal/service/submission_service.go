package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
	"github.com/blueship581/codelearn/pkg/strutil"
)

// SubmissionService 提交评测业务：提交 -> 评测 -> 统计与成就联动。
type SubmissionService struct {
	subRepo      *repository.SubmissionRepository
	problemRepo  *repository.ProblemRepository
	userRepo     *repository.UserRepository
	statRepo     *repository.UserStatRepository
	judge        *JudgeService
	achievement  *AchievementService
	logger       *slog.Logger
}

// NewSubmissionService 构造提交评测服务。
func NewSubmissionService(subRepo *repository.SubmissionRepository, problemRepo *repository.ProblemRepository,
	userRepo *repository.UserRepository, statRepo *repository.UserStatRepository, judge *JudgeService,
	achievement *AchievementService, logger *slog.Logger) *SubmissionService {
	return &SubmissionService{
		subRepo:     subRepo,
		problemRepo: problemRepo,
		userRepo:    userRepo,
		statRepo:    statRepo,
		judge:       judge,
		achievement: achievement,
		logger:      logger,
	}
}

// Submit 提交代码并同步评测：ACM 风格逐用例运行。
func (s *SubmissionService) Submit(ctx context.Context, userID primitive.ObjectID, problemID primitive.ObjectID, req *dto.SubmitRequest) (*dto.SubmissionResponse, error) {
	problem, err := s.problemRepo.FindByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, repository.ErrProblemNotFound) {
			return nil, util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if problem.Status != constants.StatusPublished {
		return nil, util.WrapAppError(constants.CodeProblemLocked, constants.MsgProblemLocked, nil)
	}
	if !strutil.Contains(problem.Languages, req.Language) {
		return nil, util.WrapAppError(constants.CodeJudgeLanguage, constants.MsgJudgeLanguage, nil)
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
	}

	sub := &model.Submission{
		UserID:       userID,
		Username:     user.Username,
		ProblemID:    problemID,
		ProblemTitle: problem.Title,
		Language:     req.Language,
		Code:         req.Code,
		Status:       constants.SubmissionJudging,
	}
	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogSubmissionCreated, "submission_id", sub.ID.Hex(), "problem_id", problemID.Hex(), "language", req.Language)

	results, status, score, runtimeMs, errMsg := s.judge.Judge(ctx, req.Language, req.Code, problem.TestCases, problem.TimeLimit)
	pointsAwarded := int64(0)
	if status == constants.SubmissionAccepted {
		pointsAwarded = int64(problem.Points)
	}
	if err := s.subRepo.UpdateResult(ctx, sub.ID, status, score, pointsAwarded, runtimeMs, results, errMsg); err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}

	// 统计联动（原子操作，无事务依赖）：题目提交数、用户统计。
	_ = s.problemRepo.IncSubmit(ctx, problemID)
	dayKey := util.SignInDailyKey(time.Now())
	_ = s.statRepo.AddSubmission(ctx, userID, req.Language, status == constants.SubmissionAccepted, dayKey)

	if status == constants.SubmissionAccepted {
		// 首次通过才累计积分/解题数/通过数，避免重复刷分。
		already, err := s.subRepo.CountAcceptedByUser(ctx, userID, problemID)
		if err == nil && already <= 1 {
			_ = s.userRepo.AddPoints(ctx, userID, pointsAwarded)
			_ = s.userRepo.MarkSolved(ctx, userID)
			_ = s.problemRepo.IncAccepted(ctx, problemID)
		}
		// 成就检查：首次通过/完成 N 题/首次通过困难题。
		s.achievement.CheckAfterSubmission(ctx, userID, status, problem)
	}

	sub.Status = status
	sub.Score = score
	sub.PointsAwarded = pointsAwarded
	sub.RuntimeMs = runtimeMs
	sub.Results = results
	sub.ErrorMessage = errMsg
	resp := dto.ToSubmissionResponse(sub)
	return &resp, nil
}

// Get 查询提交记录（本人或管理员）。
func (s *SubmissionService) Get(ctx context.Context, submissionID, userID primitive.ObjectID, role string) (*dto.SubmissionResponse, error) {
	sub, err := s.subRepo.FindByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, repository.ErrSubmissionNotFound) {
			return nil, util.WrapAppError(constants.CodeSubmissionNotFound, constants.MsgSubmissionNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if sub.UserID != userID && role != constants.RoleAdmin {
		return nil, util.WrapAppError(constants.CodeSubmissionDenied, constants.MsgSubmissionDenied, nil)
	}
	resp := dto.ToSubmissionResponse(sub)
	return &resp, nil
}

// List 分页查询提交（支持按用户/题目/状态过滤）。
func (s *SubmissionService) List(ctx context.Context, userID primitive.ObjectID, role string, page, pageSize int64, problemID string, status string) ([]dto.SubmissionResponse, int64, error) {
	filter := bson.M{}
	if problemID != "" {
		pid, err := primitive.ObjectIDFromHex(problemID)
		if err != nil {
			return nil, 0, util.WrapAppError(constants.CodeBadRequest, constants.MsgBadRequest, err)
		}
		filter["problem_id"] = pid
	}
	if status != "" {
		filter["status"] = status
	}
	// 学生只能看自己的提交；管理员可查看全部（按 problem_id/status 过滤）。
	if role != constants.RoleAdmin {
		filter["user_id"] = userID
	}
	subs, total, err := s.subRepo.List(ctx, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.SubmissionResponse, 0, len(subs))
	for _, sub := range subs {
		out = append(out, dto.ToSubmissionResponse(sub))
	}
	return out, total, nil
}
