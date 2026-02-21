package cmd

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
)

func TestCommandExit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		out     any
		wantErr error
	}{
		{name: "success", out: &bytes.Buffer{}, wantErr: ErrStop},
		{name: "write error", out: errWriter{}, wantErr: writeError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{
				Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:    out,
			}

			err := commandExit(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNewExitCommand(t *testing.T) {
	t.Parallel()

	cmd := newExitCommand()
	if cmd == nil || cmd.Callback == nil {
		t.Fatal("expected non-nil command and callback")
	}
	if cmd.name != "exit" {
		t.Fatalf("expected name exit, got %q", cmd.name)
	}
}

func TestCommandHelp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		out     any
		wantErr error
	}{
		{name: "success", out: &bytes.Buffer{}, wantErr: nil},
		{name: "write error", out: errWriter{}, wantErr: writeError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{
				Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:    out,
			}

			err := commandHelp(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNewHelpCommand(t *testing.T) {
	t.Parallel()

	cmd := newHelpCommand()
	if cmd == nil || cmd.Callback == nil {
		t.Fatal("expected non-nil command and callback")
	}
	if cmd.name != "help" {
		t.Fatalf("expected name help, got %q", cmd.name)
	}
}
