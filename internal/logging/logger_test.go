package logging

import (
	"os"
	"path/filepath"
	"testing"

	"log/slog"
)

func TestNewLoggerWritesFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	logger := NewLogger(MyHandler{Level: slog.LevelInfo})
	if logger == nil {
		t.Fatal("expected logger not nil")
	}
	logger.Info("test log")

	logPath := filepath.Join(tmpDir, "pokedex.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected log file to contain data")
	}
}
