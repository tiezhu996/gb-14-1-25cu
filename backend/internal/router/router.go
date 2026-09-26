package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/database"
	"github.com/blueship581/codelearn/internal/handler"
	"github.com/blueship581/codelearn/internal/middleware"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/service"
)

// Deps 路由装配依赖。
type Deps struct {
	DB           *database.Mongo
	Redis        *redis.Client
	Logger       *slog.Logger
	JWTSecret    string
	JudgeTimeout int
}

// New 装配全部路由（handler -> service -> repository 单向依赖、构造器注入）。
func New(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(deps.Logger), middleware.ErrorHandler(deps.Logger), middleware.CORS())

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": constants.CodeOK, "message": constants.MsgHealthOK, "data": gin.H{"status": "ok"}})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": constants.CodeOK, "message": constants.MsgOK, "data": gin.H{"service": "codelearn"}})
	})

	db := deps.DB.DB
	// repository
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	discussionRepo := repository.NewDiscussionRepository(db)
	achievementRepo := repository.NewAchievementRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	statRepo := repository.NewUserStatRepository(db)

	// service
	userService := service.NewUserService(userRepo, statRepo, deps.Logger, deps.JWTSecret)
	achievementService := service.NewAchievementService(achievementRepo, submissionRepo, userRepo, deps.Logger)
	courseService := service.NewCourseService(courseRepo, statRepo, deps.Logger)
	problemService := service.NewProblemService(problemRepo, deps.Logger)
	judgeService := service.NewJudgeService(deps.Logger, deps.JudgeTimeout)
	submissionService := service.NewSubmissionService(submissionRepo, problemRepo, userRepo, statRepo, judgeService, achievementService, deps.Logger)
	discussionService := service.NewDiscussionService(discussionRepo, userRepo, problemRepo, deps.Logger)
	leaderboardService := service.NewLeaderboardService(submissionRepo, userRepo, deps.Logger)
	dashboardService := service.NewDashboardService(statRepo, userRepo, deps.Logger)
	auditService := service.NewAuditService(auditRepo, deps.Logger)

	// handler
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	courseHandler := handler.NewCourseHandler(courseService)
	problemHandler := handler.NewProblemHandler(problemService)
	submissionHandler := handler.NewSubmissionHandler(submissionService)
	discussionHandler := handler.NewDiscussionHandler(discussionService)
	achievementHandler := handler.NewAchievementHandler(achievementService)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	auditHandler := handler.NewAuditHandler(auditService)

	auth := middleware.Auth(deps.JWTSecret)
	audit := middleware.Audit(auditService)
	admin := middleware.RequireRole(constants.RoleAdmin)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", audit, authHandler.Register)
		api.POST("/auth/login", audit, authHandler.Login)
		api.POST("/auth/sign-in", auth, middleware.RateLimit(deps.Redis, 5, time.Minute), audit, userHandler.SignIn)

		api.GET("/users/me", auth, userHandler.GetMe)
		api.PUT("/users/me", auth, audit, userHandler.UpdateProfile)
		api.GET("/users", auth, admin, userHandler.List)
		api.PUT("/users/:id/status", auth, admin, audit, userHandler.UpdateStatus)
		api.PUT("/users/:id/role", auth, admin, audit, userHandler.UpdateRole)

		registerCourseRoutes(api, courseHandler, auth, admin, audit)
		registerProblemRoutes(api, problemHandler, submissionHandler, discussionHandler, auth, admin, audit)
		registerSubmissionRoutes(api, submissionHandler, auth, admin, audit)
		registerDiscussionRoutes(api, discussionHandler, auth, admin, audit)
		registerAchievementRoutes(api, achievementHandler, auth)
		registerLeaderboardRoutes(api, leaderboardHandler, auth)
		registerDashboardRoutes(api, dashboardHandler, auth)
		registerAuditRoutes(api, auditHandler, auth, admin)
	}

	return r
}
