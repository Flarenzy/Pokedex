package cmd

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
)

const exploreFixture = `{"pokemon_encounters":[{"pokemon":{"name":"pikachu","url":"u1"}},{"pokemon":{"name":"bulbasaur","url":"u2"}}]}`

func TestCommandExplore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		body         string
		clientErr    error
		out          any
		wantErr      error
		wantContains []string
	}{
		{name: "no args", args: []string{}, out: &bytes.Buffer{}, wantErr: ErrNoAreaToExplore},
		{name: "success", args: []string{"forest"}, body: exploreFixture, out: &bytes.Buffer{}, wantContains: []string{"Exploring area:  forest", "Pokemon #1: pikachu", "Pokemon #2: bulbasaur"}},
		{name: "write error", args: []string{"forest"}, body: exploreFixture, out: errWriter{}, wantErr: writeError},
		{name: "top separator write error", args: []string{"forest"}, body: exploreFixture, out: &failOnWriteN{n: 2, err: writeError}, wantErr: writeError},
		{name: "bottom separator write error", args: []string{"forest"}, body: exploreFixture, out: &failOnWriteN{n: 5, err: writeError}, wantErr: writeError},
		{name: "final newline write error", args: []string{"forest"}, body: exploreFixture, out: &failOnWriteN{n: 6, err: writeError}, wantErr: writeError},
		{name: "pokemon fetch parse error is swallowed", args: []string{"forest"}, body: "{", out: &bytes.Buffer{}, wantContains: []string{"Exploring area:  forest", "==================================="}},
		{name: "pokemon fetch api error is swallowed", args: []string{"forest"}, clientErr: httpClientDoError, out: &bytes.Buffer{}, wantContains: []string{"Exploring area:  forest", "==================================="}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := pokecache.NewCache(20 * time.Second)
			t.Cleanup(cache.Done)

			writer, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}

			c := config.Config{
				Args:       tc.args,
				Cache:      cache,
				Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:        writer,
				HTTPClient: &stubHTTPClient{body: tc.body, err: tc.clientErr},
			}

			err := commandExplore(&c)
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

func TestGetPokemonInAreaUsesCache(t *testing.T) {
	t.Parallel()

	cache := pokecache.NewCache(20 * time.Second)
	t.Cleanup(cache.Done)

	url := internal.FirstURL + "forest"
	if err := cache.Add(url, []byte(exploreFixture)); err != nil {
		t.Fatalf("cache add failed: %v", err)
	}
	out := &bytes.Buffer{}
	c := config.Config{
		Cache:      cache,
		Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
		Out:        out,
		HTTPClient: &stubHTTPClient{err: httpClientDoError},
	}

	err := getPokemonInArea(&c, url)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "Pokemon #1: pikachu") {
		t.Fatalf("expected cached output, got %q", out.String())
	}
}

func TestGetPokemonInAreaWriteError(t *testing.T) {
	t.Parallel()

	cache := pokecache.NewCache(20 * time.Second)
	t.Cleanup(cache.Done)

	c := config.Config{
		Cache:      cache,
		Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
		Out:        errWriter{},
		HTTPClient: &stubHTTPClient{body: exploreFixture},
	}

	err := getPokemonInArea(&c, internal.FirstURL+"forest")
	if !errors.Is(err, writeError) {
		t.Fatalf("expected writeError, got %v", err)
	}
}

func TestNewExploreCommand(t *testing.T) {
	t.Parallel()

	cmd := newExploreCommand()
	if cmd == nil || cmd.Callback == nil || cmd.name != "explore" {
		t.Fatal("invalid explore command")
	}
}
