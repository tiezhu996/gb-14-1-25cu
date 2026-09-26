package util

import (
	"log/slog"
	"os"
)

// NewLogger 创建全局结构化日志器（log/slog）。
// 全站 handler/service/middleware 均引用本包 logger（屎山耦合要求）。
func NewLogger(level string) *slog.Logger {
	var lv slog.Level
	switch level {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lv})
	return slog.New(h)
}
