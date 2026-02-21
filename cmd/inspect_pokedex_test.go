package cmd

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/domain"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
)

type stubPokedex struct {
	all    []pokedex.Pokemon
	byName pokedex.Pokemon
	getErr error
}

func (s *stubPokedex) Add(p pokedex.Pokemon) {}

func (s *stubPokedex) Remove(p pokedex.Pokemon) {}

func (s *stubPokedex) GetAllPokemon() []pokedex.Pokemon {
	return s.all
}

func (s *stubPokedex) GetPokemonByName(name string) (pokedex.Pokemon, error) {
	if s.getErr != nil {
		return pokedex.Pokemon{}, s.getErr
	}
	return s.byName, nil
}

func pokemonWithStatsAndTypes() pokedex.Pokemon {
	p := pokedex.Pokemon{Name: "mew", Height: 4, Weight: 40}
	p.Stats = append(p.Stats, struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	}{BaseStat: 50, Stat: struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{Name: "hp"}})
	p.Types = append(p.Types, struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	}{Type: struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{Name: "psychic"}})
	return p
}

func TestCommandInspect(t *testing.T) {
	t.Parallel()

	seedPokemon := pokedex.Pokemon{Name: "mew", Height: 4, Weight: 40}
	unknownErr := errors.New("db error")
	tests := []struct {
		name         string
		args         []string
		pokedex      domain.Pokedexer
		out          any
		wantErr      error
		wantContains []string
	}{
		{name: "no args", args: []string{}, pokedex: pokedex.NewPokedex(), out: &bytes.Buffer{}, wantErr: ErrNoPokemonToInspect},
		{name: "pokemon not found", args: []string{"mew"}, pokedex: pokedex.NewPokedex(), out: &bytes.Buffer{}, wantContains: []string{"you have not caught that pokemon"}},
		{name: "not found write error", args: []string{"mew"}, pokedex: pokedex.NewPokedex(), out: errWriter{}, wantErr: writeError},
		{name: "unexpected pokedex error", args: []string{"mew"}, pokedex: &stubPokedex{getErr: unknownErr}, out: &bytes.Buffer{}, wantErr: unknownErr},
		{name: "success", args: []string{"mew"}, pokedex: func() domain.Pokedexer { p := pokedex.NewPokedex(); p.Add(seedPokemon); return p }(), out: &bytes.Buffer{}, wantContains: []string{"Name: mew", "Height: 4", "Weight: 40", "Stats:", "Types:"}},
		{name: "details header write error", args: []string{"mew"}, pokedex: &stubPokedex{byName: seedPokemon}, out: errWriter{}, wantErr: writeError},
		{name: "stats heading write error", args: []string{"mew"}, pokedex: &stubPokedex{byName: seedPokemon}, out: &failOnWriteN{n: 2, err: writeError}, wantErr: writeError},
		{name: "stat row write error", args: []string{"mew"}, pokedex: &stubPokedex{byName: pokemonWithStatsAndTypes()}, out: &failOnWriteN{n: 3, err: writeError}, wantErr: writeError},
		{name: "types heading write error", args: []string{"mew"}, pokedex: &stubPokedex{byName: pokemonWithStatsAndTypes()}, out: &failOnWriteN{n: 4, err: writeError}, wantErr: writeError},
		{name: "type row write error", args: []string{"mew"}, pokedex: &stubPokedex{byName: pokemonWithStatsAndTypes()}, out: &failOnWriteN{n: 5, err: writeError}, wantErr: writeError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			writer, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{Args: tc.args, Pokedex: tc.pokedex, Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}), Out: writer}

			err := commandInspect(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if buf, ok := writer.(*bytes.Buffer); ok {
				for _, expected := range tc.wantContains {
					if !strings.Contains(buf.String(), expected) {
						t.Fatalf("expected output to contain %q, got %q", expected, buf.String())
					}
				}
			}
		})
	}
}

func TestNewInspectCommand(t *testing.T) {
	t.Parallel()

	cmd := newInspectCommand()
	if cmd == nil || cmd.Callback == nil || cmd.name != "inspect" {
		t.Fatal("invalid inspect command")
	}
}

func TestCommandPokedex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		seedPokedex  bool
		out          any
		wantErr      error
		wantContains []string
	}{
		{name: "empty pokedex", out: &bytes.Buffer{}, wantErr: ErrEmptyPokedex},
		{name: "success", seedPokedex: true, out: &bytes.Buffer{}, wantContains: []string{"Your Pokedex:", "  - mew"}},
		{name: "write error", seedPokedex: true, out: errWriter{}, wantErr: writeError},
		{name: "pokemon row write error", seedPokedex: true, out: &failOnWriteN{n: 2, err: writeError}, wantErr: writeError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := pokedex.NewPokedex()
			if tc.seedPokedex {
				p.Add(pokedex.Pokemon{Name: "mew"})
			}
			writer, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{Pokedex: p, Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}), Out: writer}

			err := commandPokedex(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if buf, ok := writer.(*bytes.Buffer); ok {
				for _, expected := range tc.wantContains {
					if !strings.Contains(buf.String(), expected) {
						t.Fatalf("expected output to contain %q, got %q", expected, buf.String())
					}
				}
			}
		})
	}
}

func TestNewPokedexCommand(t *testing.T) {
	t.Parallel()

	cmd := newPokedexCommand()
	if cmd == nil || cmd.Callback == nil || cmd.name != "pokedex" {
		t.Fatal("invalid pokedex command")
	}
}
