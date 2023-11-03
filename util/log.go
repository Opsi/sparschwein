package util

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// LogConfig holds the logging configuration values
type LogConfig struct {
	// Format of the log messages (json or text)
	Format string
	// Level of the log messages (debug, info, warn, error)
	Level string
}

// AddLogFlags adds the logging flags and returns the configuration struct to
// be filled with the values from the flags.
func AddLogFlags() *LogConfig {
	logConfig := &LogConfig{}

	flag.StringVar(
		&logConfig.Format,
		"log-format",
		LookupStringEnv("LOG_FORMAT", "text"),
		"Set the logging format (json or text)")
	flag.StringVar(
		&logConfig.Level,
		"log-level",
		LookupStringEnv("LOG_LEVEL", "info"),
		"Set the logging level (debug, info, warn, error)")

	return logConfig
}

// InitSlogDefault initializes the default logger with the given configuration.
func (c LogConfig) InitSlogDefault() error {
	var leveler slog.Leveler
	switch strings.ToLower(c.Level) {
	case "debug":
		leveler = slog.LevelDebug
	case "info":
		leveler = slog.LevelInfo
	case "warn":
		leveler = slog.LevelWarn
	case "error":
		leveler = slog.LevelError
	default:
		return fmt.Errorf("unknown log level %q", c.Level)
	}

	var handler slog.Handler
	switch strings.ToLower(c.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler,
		})
	default:
		return fmt.Errorf("unknown log format %q", c.Format)
	}
	slog.SetDefault(slog.New(handler))
	return nil
}
