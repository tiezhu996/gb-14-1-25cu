package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/util"
)

// JudgeService 在线评测沙箱：在服务端子进程运行用户代码，超时 10 秒自动终止。
// 支持 python / javascript / java 三种语言（运行时镜像内置解释器与 JDK）。
type JudgeService struct {
	logger      *slog.Logger
	defaultTimeout time.Duration
	workDir     string
}

// NewJudgeService 构造评测服务。
func NewJudgeService(logger *slog.Logger, timeoutSeconds int) *JudgeService {
	if timeoutSeconds <= 0 {
		timeoutSeconds = constants.DefaultJudgeTimeout
	}
	wd, _ := os.MkdirTemp("", "codelearn-judge-*")
	return &JudgeService{logger: logger, defaultTimeout: time.Duration(timeoutSeconds) * time.Second, workDir: wd}
}

// Judge 执行一次完整评测：逐用例运行用户代码，输出状态机结果。
func (s *JudgeService) Judge(ctx context.Context, language, code string, testCases []model.TestCase, timeLimit int) ([]model.JudgeResult, string, int, int64, string) {
	if timeLimit <= 0 {
		timeLimit = constants.DefaultJudgeTimeout
	}
	results := make([]model.JudgeResult, 0, len(testCases))
	passed := 0
	var totalRuntime int64
	for i, tc := range testCases {
		actual, runtimeMs, runErr := s.runCode(ctx, language, code, tc.Input, timeLimit)
		totalRuntime += runtimeMs
		res := model.JudgeResult{
			TestCaseIndex: i,
			Input:         tc.Input,
			Expected:      tc.Output,
			Actual:        actual,
		}
		if runErr != nil {
			res.ErrorMessage = runErr.Error()
			if isTimeoutErr(runErr) {
				res.Passed = false
				s.logger.Warn(constants.LogSubmissionTimeout, "test_case", i, "error", runErr.Error())
			} else {
				res.Passed = false
				s.logger.Warn(constants.LogJudgeRunFailed, "test_case", i, "error", runErr.Error())
			}
			results = append(results, res)
			// 运行错误/超时立即终止后续用例
			status := constants.SubmissionRuntimeError
			errMsg := runErr.Error()
			if isTimeoutErr(runErr) {
				status = constants.SubmissionTimeout
				errMsg = fmt.Sprintf(constants.MsgJudgeTimeout, timeLimit)
			}
			return results, status, calcScore(passed, len(testCases)), totalRuntime, errMsg
		}
		res.Passed = util.NormalizeOutput(actual) == util.NormalizeOutput(tc.Output)
		if res.Passed {
			passed++
		}
		results = append(results, res)
	}
	if passed == len(testCases) {
		s.logger.Info(constants.LogSubmissionAccepted, "passed", passed, "total", len(testCases))
		return results, constants.SubmissionAccepted, 100, totalRuntime, ""
	}
	if passed == 0 {
		return results, constants.SubmissionRuntimeError, 0, totalRuntime, "全部测试用例未通过"
	}
	s.logger.Info(constants.LogSubmissionPartial, "passed", passed, "total", len(testCases))
	return results, constants.SubmissionPartial, calcScore(passed, len(testCases)), totalRuntime, ""
}

// runCode 在独立工作目录中运行用户代码。
func (s *JudgeService) runCode(ctx context.Context, language, code, input string, timeoutSeconds int) (string, int64, error) {
	dir, err := os.MkdirTemp(s.workDir, "run-*")
	if err != nil {
		return "", 0, fmt.Errorf("创建沙箱目录失败: %w", err)
	}
	defer os.RemoveAll(dir)

	timeout := time.Duration(timeoutSeconds) * time.Second
	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	switch language {
	case constants.LanguagePython:
		file := filepath.Join(dir, "main.py")
		if err := os.WriteFile(file, []byte(code), 0o644); err != nil {
			return "", 0, fmt.Errorf("写入代码失败: %w", err)
		}
		return s.exec(ctx2, dir, input, "python3", "main.py")
	case constants.LanguageJavaScript:
		file := filepath.Join(dir, "main.js")
		if err := os.WriteFile(file, []byte(code), 0o644); err != nil {
			return "", 0, fmt.Errorf("写入代码失败: %w", err)
		}
		return s.exec(ctx2, dir, input, "node", "main.js")
	case constants.LanguageJava:
		file := filepath.Join(dir, "Main.java")
		if err := os.WriteFile(file, []byte(code), 0o644); err != nil {
			return "", 0, fmt.Errorf("写入代码失败: %w", err)
		}
		compile := exec.CommandContext(ctx2, "javac", "Main.java")
		compile.Dir = dir
		var compileErr bytes.Buffer
		compile.Stderr = &compileErr
		if err := compile.Run(); err != nil {
			return "", 0, fmt.Errorf("编译错误: %s", strings.TrimSpace(compileErr.String()))
		}
		return s.exec(ctx2, dir, input, "java", "Main")
	default:
		return "", 0, fmt.Errorf(constants.MsgJudgeLanguage, language)
	}
}

// exec 运行子进程并捕获 stdout/stderr，返回输出与耗时。
func (s *JudgeService) exec(ctx context.Context, dir, input, name string, args ...string) (string, int64, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	start := time.Now()
	err := cmd.Run()
	runtimeMs := time.Since(start).Milliseconds()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", runtimeMs, fmt.Errorf("运行超时（超过 %d 秒）", constants.DefaultJudgeTimeout)
		}
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return strings.TrimSpace(stdout.String()), runtimeMs, fmt.Errorf("运行错误: %s", truncate(msg, 500))
	}
	return stdout.String(), runtimeMs, nil
}

// isTimeoutErr 判断是否超时错误。
func isTimeoutErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "超时")
}

// calcScore 计算得分百分比。
func calcScore(passed, total int) int {
	if total == 0 {
		return 0
	}
	return passed * 100 / total
}

// truncate 截断超长错误信息。
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
