package cmd

import (
	"fmt"
	"sort"

	"github.com/Flarenzy/Pokedex/internal/config"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(c *config.Config) error
}

const helpTextBase = `
Welcome to the Pokedex!
Usage:

`

var helpText = helpTextBase

func NewCommands() map[string]*CliCommand {
	commands := make(map[string]*CliCommand)
	commands["help"] = newHelpCommand()
	commands["exit"] = newExitCommand()
	commands["map"] = newMapCommand()
	commands["mapb"] = newMapbCommand()
	commands["explore"] = newExploreCommand()
	commands["catch"] = newCatchCommand()
	commands["inspect"] = newInspectCommand()
	commands["pokedex"] = newPokedexCommand()

	keys := make([]string, 0, len(commands))
	for key := range commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	helpText = helpTextBase
	for _, key := range keys {
		helpText += fmt.Sprintf("%s: %s\n", key, commands[key].description)
	}
	return commands
}
