package cmd

import (
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"log/slog"
)

const firstURL = "https://pokeapi.co/api/v2/location-area/"

type CliCommand struct {
	name        string
	description string
	Callback    func(config *Config) error
}

type Config struct {
	Next     string
	Previous string
	Args     []string
	cache    *pokecache.Cache
	Logger   *slog.Logger
}

func NewConfig(cache *pokecache.Cache, logger *slog.Logger) *Config {
	return &Config{
		Next:     firstURL,
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
	commands["map"] = newMapCommand()
	commands["mapb"] = newMapbCommand()
	commands["explore"] = newExploreCommand()
	for k, v := range commands {
		helpText += fmt.Sprintf("%s: %s\n", k, v.description)
	}
	return commands
}
