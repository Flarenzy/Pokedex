package config

import (
	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
	"log/slog"
)

type Config struct {
	Next     string
	Previous string
	Args     []string
	Pokedex  *pokedex.Pokedex
	Cache    *pokecache.Cache
	Logger   *slog.Logger
}

func NewConfig(cache *pokecache.Cache, logger *slog.Logger, p *pokedex.Pokedex) *Config {
	return &Config{
		Next:     internal.FirstURL,
		Previous: "",
		Args:     []string{},
		Cache:    cache,
		Logger:   logger,
		Pokedex:  p,
	}
}
