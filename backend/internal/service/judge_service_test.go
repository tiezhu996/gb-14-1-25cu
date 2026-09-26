package service

import (
	"context"
	"os/exec"
	"testing"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/util"
)

func newTestJudge(t *testing.T) *JudgeService {
	t.Helper()
	return NewJudgeService(util.NewLogger("error"), constants.DefaultJudgeTimeout)
}

func hasPython() bool {
	_, err := exec.LookPath("python3")
	return err == nil
}

func TestJudgePythonAccepted(t *testing.T) {
	if !hasPython() {
		t.Skip("python3 not available")
	}
	j := newTestJudge(t)
	code := "a, b = map(int, input().split())\nprint(a + b)"
	tcs := []model.TestCase{
		{Input: "1 2", Output: "3"},
		{Input: "10 20", Output: "30"},
	}
	results, status, score, _, _ := j.Judge(context.Background(), constants.LanguagePython, code, tcs, 10)
	if status != constants.SubmissionAccepted {
		t.Errorf("status = %q, want accepted; results=%+v", status, results)
	}
	if score != 100 {
		t.Errorf("score = %d, want 100", score)
	}
	if len(results) != 2 || !results[0].Passed || !results[1].Passed {
		t.Errorf("results not all passed: %+v", results)
	}
}

func TestJudgePythonPartial(t *testing.T) {
	if !hasPython() {
		t.Skip("python3 not available")
	}
	j := newTestJudge(t)
	code := "a, b = map(int, input().split())\nprint(a + b)"
	tcs := []model.TestCase{
		{Input: "1 2", Output: "3"},
		{Input: "10 20", Output: "999"},
	}
	_, status, score, _, _ := j.Judge(context.Background(), constants.LanguagePython, code, tcs, 10)
	if status != constants.SubmissionPartial {
		t.Errorf("status = %q, want partial", status)
	}
	if score != 50 {
		t.Errorf("score = %d, want 50", score)
	}
}

func TestJudgePythonRuntimeError(t *testing.T) {
	if !hasPython() {
		t.Skip("python3 not available")
	}
	j := newTestJudge(t)
	code := "raise ValueError('boom')"
	tcs := []model.TestCase{{Input: "", Output: "3"}}
	_, status, _, _, errMsg := j.Judge(context.Background(), constants.LanguagePython, code, tcs, 10)
	if status != constants.SubmissionRuntimeError {
		t.Errorf("status = %q, want runtime_error", status)
	}
	if errMsg == "" {
		t.Error("expected error message")
	}
}

func TestJudgePythonTimeout(t *testing.T) {
	if !hasPython() {
		t.Skip("python3 not available")
	}
	j := newTestJudge(t)
	code := "while True:\n    pass"
	tcs := []model.TestCase{{Input: "", Output: "3"}}
	_, status, _, _, _ := j.Judge(context.Background(), constants.LanguagePython, code, tcs, 1)
	if status != constants.SubmissionTimeout {
		t.Errorf("status = %q, want timeout", status)
	}
}

func TestCalcScore(t *testing.T) {
	tests := []struct {
		name   string
		passed int
		total  int
		want   int
	}{
		{name: "all pass", passed: 3, total: 3, want: 100},
		{name: "half", passed: 1, total: 2, want: 50},
		{name: "zero total", passed: 0, total: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcScore(tt.passed, tt.total); got != tt.want {
				t.Errorf("calcScore(%d,%d) = %d, want %d", tt.passed, tt.total, got, tt.want)
			}
		})
	}
}
