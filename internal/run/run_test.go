package run

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/Flarenzy/Pokedex/cmd"
	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
	"github.com/chzyer/readline"
)

type scriptLine struct {
	line string
	err  error
}

type scriptReader struct {
	items []scriptLine
	idx   int
}

func (s *scriptReader) Readline() (string, error) {
	if s.idx >= len(s.items) {
		return "", io.EOF
	}
	item := s.items[s.idx]
	s.idx++
	return item.line, item.err
}

func (s *scriptReader) Close() error { return nil }

type runCache struct{ done bool }

func (r *runCache) Get(key string) ([]byte, error)   { return nil, nil }
func (r *runCache) Add(key string, val []byte) error { return nil }
func (r *runCache) Done()                            { r.done = true }

func testConfig(cache *runCache) *config.Config {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return &config.Config{
		Cache:   cache,
		Logger:  logger,
		Pokedex: pokedex.NewPokedex(),
		Out:     io.Discard,
	}
}

func TestCleanInput(t *testing.T) {
	t.Parallel()

	got := cleanInput("  HeLp   MAPB  FoO ")
	want := []string{"help", "mapb", "foo"}
	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected token %d=%q, got %q", i, want[i], got[i])
		}
	}
}

func TestRunExecutesCommandAndStops(t *testing.T) {
	t.Parallel()

	cache := &runCache{}
	c := testConfig(cache)

	called := false
	commands := map[string]*cmd.CliCommand{
		"catch": {Callback: func(cfg *config.Config) error {
			called = true
			if len(cfg.Args) != 2 || cfg.Args[0] != "mew" || cfg.Args[1] != "mewtwo" {
				t.Fatalf("args not parsed as expected: %v", cfg.Args)
			}
			return cmd.ErrStop
		}},
	}

	in := &scriptReader{items: []scriptLine{{line: "CATCH   mew   MEWTWO", err: nil}}}
	err := Run(c, in, commands)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !called {
		t.Fatal("expected command callback to be called")
	}
	if !cache.done {
		t.Fatal("expected cache.Done to be called on ErrStop")
	}
}

func TestRunReturnsCommandError(t *testing.T) {
	t.Parallel()

	cache := &runCache{}
	c := testConfig(cache)
	expected := errors.New("boom")
	commands := map[string]*cmd.CliCommand{
		"x": {Callback: func(cfg *config.Config) error { return expected }},
	}

	in := &scriptReader{items: []scriptLine{{line: "x", err: nil}}}
	err := Run(c, in, commands)
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	if cache.done {
		t.Fatal("did not expect cache.Done on non-stop error")
	}
}

func TestRunContinuesOnUnknownCommand(t *testing.T) {
	t.Parallel()

	cache := &runCache{}
	c := testConfig(cache)
	called := false
	commands := map[string]*cmd.CliCommand{
		"exit": {Callback: func(cfg *config.Config) error { called = true; return cmd.ErrStop }},
	}
	in := &scriptReader{items: []scriptLine{{line: "unknown", err: nil}, {line: "exit", err: nil}}}

	err := Run(c, in, commands)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !called {
		t.Fatal("expected known command callback to run")
	}
}

func TestRunHandlesReadlineErrorAndProcessesLine(t *testing.T) {
	t.Parallel()

	cache := &runCache{}
	c := testConfig(cache)
	called := false
	commands := map[string]*cmd.CliCommand{
		"stop": {Callback: func(cfg *config.Config) error { called = true; return cmd.ErrStop }},
	}
	in := &scriptReader{items: []scriptLine{{line: "stop", err: errors.New("transient readline error")}}}

	err := Run(c, in, commands)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !called {
		t.Fatal("expected callback to run even when readline returned non-interrupt error")
	}
}

func TestNewReadlineInputAndMethods(t *testing.T) {
	cfg := &readline.Config{
		Prompt:              ">",
		Stdin:               io.NopCloser(strings.NewReader("hello\n")),
		Stdout:              io.Discard,
		Stderr:              io.Discard,
		ForceUseInteractive: false,
	}
	rl, err := readline.NewEx(cfg)
	if err != nil {
		t.Fatalf("failed to create readline instance: %v", err)
	}

	in := NewReadlineInput(rl)
	if in == nil || in.rl != rl {
		t.Fatal("expected NewReadlineInput to wrap instance")
	}
	_, _ = in.Readline()
	if err := in.Close(); err != nil {
		t.Fatalf("expected close nil, got %v", err)
	}
}
