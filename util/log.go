package util

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// InitSlogDefault initializes the default logger with the LOG_FORMAT and
// LOG_LEVEL environment variables.
func InitSlogDefault() error {
	var (
		format = LookupStringEnv("LOG_FORMAT", "text")
		level  = LookupStringEnv("LOG_LEVEL", "info")
	)
	var leveler slog.Leveler
	switch strings.ToLower(level) {
	case "debug":
		leveler = slog.LevelDebug
	case "info":
		leveler = slog.LevelInfo
	case "warn":
		leveler = slog.LevelWarn
	case "error":
		leveler = slog.LevelError
	default:
		return fmt.Errorf("unknown log level %q", level)
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler,
		})
	default:
		return fmt.Errorf("unknown log format %q", format)
	}
	slog.SetDefault(slog.New(handler))
	return nil
}
