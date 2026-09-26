package seed

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// Seed 初始化演示数据：管理员、学生、课程、题目、成就定义。
// 全部幂等：已存在则跳过。
func Seed(ctx context.Context, db *mongo.Database, logger *slog.Logger) error {
	if err := repository.NewAchievementRepository(db).Seed(ctx, constants.AchievementDefs); err != nil {
		return err
	}

	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	statRepo := repository.NewUserStatRepository(db)

	// 管理员
	adminHash, _ := util.HashPassword("admin123")
	admin := &model.User{
		Username:     "admin",
		Email:        "admin@codelearn.dev",
		PasswordHash: adminHash,
		Nickname:     "管理员",
		Role:         constants.RoleAdmin,
		Status:       constants.UserStatusActive,
		Points:       1200,
		SolvedCount:  38,
	}
	if _, err := userRepo.FindByUsername(ctx, "admin"); err != nil {
		if err := userRepo.Create(ctx, admin); err != nil {
			return err
		}
		_ = statRepo.Upsert(ctx, admin.ID)
	}

	// 学生
	studentHash, _ := util.HashPassword("student123")
	student := &model.User{
		Username:     "student",
		Email:        "student@codelearn.dev",
		PasswordHash: studentHash,
		Nickname:     "小码同学",
		Role:         constants.RoleStudent,
		Status:       constants.UserStatusActive,
		Points:       320,
		SolvedCount:  12,
	}
	if _, err := userRepo.FindByUsername(ctx, "student"); err != nil {
		if err := userRepo.Create(ctx, student); err != nil {
			return err
		}
		_ = statRepo.Upsert(ctx, student.ID)
		_ = statRepo.AddSubmission(ctx, student.ID, constants.LanguagePython, true, util.SignInDailyKey(time.Now()))
	}

	// 课程
	if _, err := courseRepo.FindByTitle(ctx, "Python 入门到进阶"); err != nil {
		course := &model.Course{
			Title:       "Python 入门到进阶",
			Description: "从零开始学习 Python 语法、数据结构与常用库，配套实战练习。",
			Cover:       "https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=800",
			Markdown:    "# Python 入门到进阶\n\n本课程覆盖 Python 基础语法、函数、面向对象、异常处理与常用标准库。\n\n```python\nprint(\"Hello, CodeLearn!\")\n```\n",
			Chapters: []model.Chapter{
				{Title: "变量与数据类型", Content: "## 变量与数据类型\n\nPython 是动态类型语言，常见类型包括 int、float、str、list、dict。", Duration: 20},
				{Title: "条件与循环", Content: "## 条件与循环\n\n使用 if/elif/else 与 for/while 控制流程。", Duration: 30},
				{Title: "函数与模块", Content: "## 函数与模块\n\n使用 def 定义函数，import 引入模块。", Duration: 25},
			},
			Difficulty: constants.DifficultyEasy,
			Status:     constants.StatusPublished,
			AuthorID:   admin.ID,
			AuthorName: admin.Username,
		}
		if err := courseRepo.Create(ctx, course); err != nil {
			return err
		}
	}

	if _, err := courseRepo.FindByTitle(ctx, "Web 开发实战"); err != nil {
		course := &model.Course{
			Title:       "Web 开发实战",
			Description: "使用 Flask/FastAPI 构建 REST API，掌握前后端分离开发流程。",
			Cover:       "https://images.unsplash.com/photo-1547658719-da2b51169166?w=800",
			Markdown:    "# Web 开发实战\n\n学习路由、请求参数、JSON 响应与数据库集成。\n",
			Chapters: []model.Chapter{
				{Title: "HTTP 与 REST", Content: "## HTTP 与 REST\n\n理解 GET/POST/PUT/DELETE 与资源建模。", Duration: 30},
				{Title: "构建第一个 API", Content: "## 构建第一个 API\n\n使用 FastAPI 创建 /hello 接口。", Duration: 40},
			},
			Difficulty: constants.DifficultyMedium,
			Status:     constants.StatusPublished,
			AuthorID:   admin.ID,
			AuthorName: admin.Username,
		}
		if err := courseRepo.Create(ctx, course); err != nil {
			return err
		}
	}

	// 题目
	seedProblem(ctx, problemRepo, &model.Problem{
		Title:        "两数之和",
		Description:  "给定两个整数 a 和 b，输出它们的和。\n\n**输入**：一行两个整数 a b。\n**输出**：一个整数，表示 a+b。\n\n示例：\n```\n输入：1 2\n输出：3\n```\n",
		Difficulty:   constants.DifficultyEasy,
		Languages:    []string{constants.LanguagePython, constants.LanguageJavaScript, constants.LanguageJava},
		Tags:         []string{"入门", "数学"},
		TestCases:    []model.TestCase{{Input: "1 2", Output: "3"}, {Input: "10 20", Output: "30"}, {Input: "-5 8", Output: "3"}},
		TimeLimit:    10,
		Status:       constants.StatusPublished,
		Points:       constants.DifficultyPoints[constants.DifficultyEasy],
		CreatedBy:    admin.ID,
		CreatedByName: admin.Username,
	})
	seedProblem(ctx, problemRepo, &model.Problem{
		Title:        "判断奇偶",
		Description:  "给定一个整数 n，判断它是奇数还是偶数。若为偶数输出 even，否则输出 odd。\n\n示例：\n```\n输入：4\n输出：even\n```\n",
		Difficulty:   constants.DifficultyEasy,
		Languages:    []string{constants.LanguagePython, constants.LanguageJavaScript, constants.LanguageJava},
		Tags:         []string{"入门", "条件"},
		TestCases:    []model.TestCase{{Input: "4", Output: "even"}, {Input: "7", Output: "odd"}, {Input: "0", Output: "even"}},
		TimeLimit:    10,
		Status:       constants.StatusPublished,
		Points:       constants.DifficultyPoints[constants.DifficultyEasy],
		CreatedBy:    admin.ID,
		CreatedByName: admin.Username,
	})
	seedProblem(ctx, problemRepo, &model.Problem{
		Title:        "斐波那契数列",
		Description:  "给定 n，输出斐波那契数列第 n 项（F(1)=1, F(2)=1, F(n)=F(n-1)+F(n-2)）。\n\n示例：\n```\n输入：6\n输出：8\n```\n",
		Difficulty:   constants.DifficultyMedium,
		Languages:    []string{constants.LanguagePython, constants.LanguageJavaScript, constants.LanguageJava},
		Tags:         []string{"递推", "动态规划"},
		TestCases:    []model.TestCase{{Input: "1", Output: "1"}, {Input: "6", Output: "8"}, {Input: "10", Output: "55"}},
		TimeLimit:    10,
		Status:       constants.StatusPublished,
		Points:       constants.DifficultyPoints[constants.DifficultyMedium],
		CreatedBy:    admin.ID,
		CreatedByName: admin.Username,
	})
	seedProblem(ctx, problemRepo, &model.Problem{
		Title:        "最长公共前缀",
		Description:  "给定多行字符串，第一行为 n（字符串数量），随后 n 行字符串，输出它们的最长公共前缀。\n\n示例：\n```\n输入：3\\nflower\\nflow\\nflight\n输出：fl\n```\n",
		Difficulty:   constants.DifficultyHard,
		Languages:    []string{constants.LanguagePython, constants.LanguageJavaScript, constants.LanguageJava},
		Tags:         []string{"字符串"},
		TestCases:    []model.TestCase{{Input: "3\nflower\nflow\nflight", Output: "fl"}, {Input: "2\ndog\nracecar", Output: ""}},
		TimeLimit:    10,
		Status:       constants.StatusPublished,
		Points:       constants.DifficultyPoints[constants.DifficultyHard],
		CreatedBy:    admin.ID,
		CreatedByName: admin.Username,
	})

	logger.Info("seed data ready")
	return nil
}

// seedProblem 幂等创建题目。
func seedProblem(ctx context.Context, repo *repository.ProblemRepository, p *model.Problem) {
	if _, err := repo.FindByTitle(ctx, p.Title); err == nil {
		return
	}
	_ = repo.Create(ctx, p)
}
