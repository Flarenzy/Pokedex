package config

import (
	"io"
	"log/slog"
	"os"

	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/http"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
)

type Config struct {
	Next        string
	Previous    string
	Args        []string
	Pokedex     *pokedex.Pokedex
	Cache       *pokecache.Cache
	Logger      *slog.Logger
	Out         io.Writer
	HTTPClient  http.HTTPClienter
	RandFloat64 func() float64
}

func NewConfig(
	cache *pokecache.Cache,
	logger *slog.Logger,
	p *pokedex.Pokedex,
	client http.HTTPClienter,
	randFloat64 func() float64,
) *Config {
	return &Config{
		Next:        internal.FirstURL,
		Previous:    "",
		Args:        []string{},
		Cache:       cache,
		Logger:      logger,
		Pokedex:     p,
		Out:         os.Stdout,
		HTTPClient:  client,
		RandFloat64: randFloat64,
	}
}
