package cmd

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
)

func TestNoPokemonToCatch(t *testing.T) {
	t.Parallel()
	c := config.Config{
		Next:     "",
		Previous: "",
		Args:     []string{},
		Pokedex:  nil,
		Cache:    nil,
		Logger: logging.NewLogger(logging.MyHandler{
			Level: slog.LevelError,
		}),
		Out:         nil,
		HTTPClient:  nil,
		RandFloat64: nil,
	}

	err := commandCatch(&c)
	if !errors.Is(err, ErrNoPokemon) {
		t.Error("Expected ErrNoPokemon", ErrNoPokemon)
	}
}
