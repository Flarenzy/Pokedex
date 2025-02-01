package cmd

import (
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"log/slog"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(config *Config) error
}

type Config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
	Logger   *slog.Logger
}

func NewConfig(cache *pokecache.Cache, logger *slog.Logger) *Config {
	return &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
		Logger:   logger,
	}
}

var helpText = `
Welcome to the Pokedex!
Usage:

`

func NewCommands() map[string]*CliCommand {
	commands := make(map[string]*CliCommand)
	commands["help"] = newHelpCommand()
	commands["exit"] = newExitCommand()
	commands["map"] = NewMapCommand()
	commands["mapb"] = NewMapbCommand()
	for k, v := range commands {
		helpText += fmt.Sprintf("%s: %s\n", k, v.description)
	}
	return commands
}
