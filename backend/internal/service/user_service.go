package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// UserService 用户业务：注册/登录/签到/资料管理。
type UserService struct {
	repo   *repository.UserRepository
	stat   *repository.UserStatRepository
	logger *slog.Logger
	secret string
}

// NewUserService 构造用户服务。
func NewUserService(repo *repository.UserRepository, stat *repository.UserStatRepository, logger *slog.Logger, secret string) *UserService {
	return &UserService{repo: repo, stat: stat, logger: logger, secret: secret}
}

// Register 注册：校验唯一性 -> 哈希密码 -> 落库 -> 初始化统计。
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	if _, err := s.repo.FindByUsername(ctx, req.Username); err == nil {
		return nil, util.WrapAppError(constants.CodeUserExists, constants.MsgUserExists, nil)
	}
	if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, util.WrapAppError(constants.CodeUserExists, constants.MsgUserExists, nil)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error(constants.LogUserRegisterFailed, "username", req.Username, "error", err.Error())
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Nickname:     nickname,
		Role:         constants.RoleStudent,
		Status:       constants.UserStatusActive,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, util.WrapAppError(constants.CodeUserExists, constants.MsgUserExists, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if err := s.stat.Upsert(ctx, user.ID); err != nil {
		s.logger.Error(constants.LogUserRegisterFailed, "username", req.Username, "error", err.Error())
	}
	token, err := util.GenerateToken(s.secret, user.ID, user.Username, user.Role, user.Status)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserRegisterSuccess, "username", user.Username, "user_id", user.ID.Hex())
	return &dto.LoginResponse{Token: token, User: dto.ToUserResponse(user)}, nil
}

// Login 登录：校验账号状态 -> 校验密码 -> 签发 JWT。
func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if errors.Is(err, repository.ErrUserNotFound) {
		s.logger.Warn(constants.LogUserLoginFailed, "username", req.Username, "reason", "user not found")
		return nil, util.WrapAppError(constants.CodePasswordIncorrect, constants.MsgUserWrongPass, err)
	}
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if user.Status == constants.UserStatusBanned {
		return nil, util.WrapAppError(constants.CodeUserBanned, constants.MsgUserBanned, nil)
	}
	if !util.CheckPassword(user.PasswordHash, req.Password) {
		s.logger.Warn(constants.LogUserLoginFailed, "username", req.Username, "reason", "wrong password")
		return nil, util.WrapAppError(constants.CodePasswordIncorrect, constants.MsgUserWrongPass, nil)
	}
	token, err := util.GenerateToken(s.secret, user.ID, user.Username, user.Role, user.Status)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserLoginSuccess, "username", user.Username, "user_id", user.ID.Hex())
	return &dto.LoginResponse{Token: token, User: dto.ToUserResponse(user)}, nil
}

// GetProfile 获取个人资料。
func (s *UserService) GetProfile(ctx context.Context, userID primitive.ObjectID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}

// UpdateProfile 更新个人资料。
func (s *UserService) UpdateProfile(ctx context.Context, userID primitive.ObjectID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	update := bson.M{}
	if req.Nickname != "" {
		update["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		update["avatar"] = req.Avatar
	}
	if req.Email != "" {
		if existing, err := s.repo.FindByEmail(ctx, req.Email); err == nil && existing.ID != userID {
			return nil, util.WrapAppError(constants.CodeUserExists, constants.MsgUserExists, nil)
		}
		update["email"] = req.Email
	}
	if err := s.repo.UpdateProfile(ctx, userID, update); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserProfileUpdated, "user_id", userID.Hex())
	return s.GetProfile(ctx, userID)
}

// SignIn 每日签到：连续签到天数状态机（昨天签过则 +1，否则重置为 1；今天已签则报错）。
func (s *UserService) SignIn(ctx context.Context, userID primitive.ObjectID) (*dto.SignInResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	now := time.Now()
	today := util.SignInDailyKey(now)
	if user.LastSignInAt != nil && util.SignInDailyKey(*user.LastSignInAt) == today {
		return nil, util.WrapAppError(constants.CodeUserSignInLock, constants.MsgUserSignedToday, nil)
	}
	streak := 1
	if user.LastSignInAt != nil {
		yesterday := now.AddDate(0, 0, -1)
		if util.SignInDailyKey(*user.LastSignInAt) == util.SignInDailyKey(yesterday) {
			streak = user.StreakDays + 1
		}
	}
	if err := s.repo.UpdateSignIn(ctx, userID, streak, now); err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserSignInSuccess, "user_id", userID.Hex(), "streak", streak)
	return &dto.SignInResponse{StreakDays: streak, FirstToday: true}, nil
}

// List 分页查询用户（管理员）。
func (s *UserService) List(ctx context.Context, page, pageSize int64, role, status string) ([]dto.UserResponse, int64, error) {
	filter := bson.M{}
	if role != "" {
		filter["role"] = role
	}
	if status != "" {
		filter["status"] = status
	}
	users, total, err := s.repo.List(ctx, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error(constants.LogUserListFailed, "error", err.Error())
		return nil, 0, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.ToUserResponse(u))
	}
	return out, total, nil
}

// UpdateStatus 管理员更新用户状态。
func (s *UserService) UpdateStatus(ctx context.Context, adminID, userID primitive.ObjectID, status string) error {
	if adminID == userID {
		return util.WrapAppError(constants.CodeUserRoleChange, "不能修改自己的账号状态", nil)
	}
	if err := s.repo.UpdateStatus(ctx, userID, status); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogUserSignInSuccess, "admin_id", adminID.Hex(), "user_id", userID.Hex(), "status", status)
	return nil
}

// UpdateRole 管理员更新用户角色（RBAC）。
func (s *UserService) UpdateRole(ctx context.Context, userID primitive.ObjectID, role string) error {
	if !constants.ValidRole(role) {
		return util.WrapAppError(constants.CodeValidation, fmt.Sprintf("角色 %s 不合法", role), nil)
	}
	if err := s.repo.UpdateRole(ctx, userID, role); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	return nil
}
