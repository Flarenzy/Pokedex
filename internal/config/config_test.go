package config

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
)

type stubCache struct{}

func (s stubCache) Get(key string) ([]byte, error)   { return nil, nil }
func (s stubCache) Add(key string, val []byte) error { return nil }
func (s stubCache) Done()                            {}

type stubPokedex struct{}

func (s stubPokedex) Add(p pokedex.Pokemon)            {}
func (s stubPokedex) Remove(p pokedex.Pokemon)         {}
func (s stubPokedex) GetAllPokemon() []pokedex.Pokemon { return nil }
func (s stubPokedex) GetPokemonByName(name string) (pokedex.Pokemon, error) {
	return pokedex.Pokemon{}, nil
}

type stubHTTPClient struct{}

func (s stubHTTPClient) Do(req *http.Request) (*http.Response, error) { return nil, nil }

func TestNewConfig(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	randFn := func() float64 { return 0.5 }

	cache := stubCache{}
	p := stubPokedex{}
	client := stubHTTPClient{}

	c := NewConfig(cache, logger, p, client, randFn)
	if c == nil {
		t.Fatal("expected config not nil")
	}
	if c.Next != internal.FirstURL || c.AreaURL != internal.FirstURL {
		t.Fatalf("expected default first URL, got next=%q area=%q", c.Next, c.AreaURL)
	}
	if c.PokemonURL != internal.SecondURL {
		t.Fatalf("expected default second URL, got %q", c.PokemonURL)
	}
	if c.Previous != "" {
		t.Fatalf("expected empty previous, got %q", c.Previous)
	}
	if len(c.Args) != 0 {
		t.Fatalf("expected empty args, got %v", c.Args)
	}
	if c.Logger != logger {
		t.Fatal("expected logger to be assigned")
	}
	if c.Cache == nil || c.Pokedex == nil || c.HTTPClient == nil {
		t.Fatal("expected dependencies to be assigned")
	}
	if c.Out != os.Stdout {
		t.Fatal("expected default output to os.Stdout")
	}
	if c.RandFloat64 == nil || c.RandFloat64() != 0.5 {
		t.Fatal("expected RandFloat64 to be assigned")
	}
}
