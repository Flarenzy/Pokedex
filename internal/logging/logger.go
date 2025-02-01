package logging

import (
	"log/slog"
	"os"
)

type MyHandler struct {
	Level slog.Level
}

func NewLogger(level MyHandler) *slog.Logger {
	f, err := os.OpenFile("pokedex.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
	}
	logger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: level.Level}))
	return logger
}
