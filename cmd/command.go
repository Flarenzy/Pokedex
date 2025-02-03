package cmd

import (
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/config"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(c *config.Config) error
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
	commands["catch"] = newCatchCommand()
	commands["inspect"] = newInspectCommand()
	commands["pokedex"] = newPokedexCommand()
	for k, v := range commands {
		helpText += fmt.Sprintf("%s: %s\n", k, v.description)
	}
	return commands
}
