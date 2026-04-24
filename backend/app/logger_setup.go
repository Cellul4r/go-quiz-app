package main

import (
	"log/slog"
	"os"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
)

func newAppLogger(cfg config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.AppEnv == "development" || cfg.AppDebug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
