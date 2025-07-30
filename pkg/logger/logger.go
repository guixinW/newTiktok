package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New 创建并返回一个配置好的 slog.Logger 实例
// level 参数可以是 "debug", "info", "warn", "error"
func New(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	// 使用 JSON Handler，方便后续的日志收集和分析
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
		// 如果需要，可以添加 Source，以便追踪日志来源文件和行号
		// AddSource: true,
	})

	return slog.New(handler)
}
