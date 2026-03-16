package infrastructure

import (
	"log/slog"
	"os"
)

// NewLogger creates a new structured logger
func NewLogger(config *Config) *slog.Logger {
	var handler slog.Handler

	// Parse log level
	level := parseLogLevel(config.LogLevel)

	// Create handler based on format
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch config.LogFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// parseLogLevel parses the log level string
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Made with Bob
