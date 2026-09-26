package service

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// AuditService 操作审计日志业务。
type AuditService struct {
	repo   *repository.AuditRepository
	logger *slog.Logger
}

// NewAuditService 构造审计服务。
func NewAuditService(repo *repository.AuditRepository, logger *slog.Logger) *AuditService {
	return &AuditService{repo: repo, logger: logger}
}

// Create 写入审计日志（middleware 与 service 埋点复用）。
func (s *AuditService) Create(ctx context.Context, log *model.AuditLog) {
	if err := s.repo.Create(ctx, log); err != nil {
		s.logger.Error(constants.LogAuditWriteFailed, "error", err.Error(), "path", log.Path)
		return
	}
	s.logger.Info(constants.LogAuditWritten, "username", log.Username, "action", log.Action, "path", log.Path)
}

// List 分页查询审计日志（管理员）。
func (s *AuditService) List(ctx context.Context, page, pageSize int64, username, entity string) ([]dto.AuditLogResponse, int64, error) {
	filter := bson.M{}
	if username != "" {
		filter["username"] = username
	}
	if entity != "" {
		filter["entity"] = entity
	}
	list, total, err := s.repo.List(ctx, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.AuditLogResponse, 0, len(list))
	for _, l := range list {
		out = append(out, dto.ToAuditLogResponse(l))
	}
	return out, total, nil
}
