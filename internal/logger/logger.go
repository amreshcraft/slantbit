package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func InitLogger(env string) (*slog.Logger, func()) {
	// 1. Ensure 'logs' directory exists
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	logPath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	cleanup := func() {
		_ = logFile.Close()
	}

	return logger, cleanup
}