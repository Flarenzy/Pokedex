package cmd

import "fmt"

type CliCommand struct {
	name        string
	description string
	Callback    func() error
}

var helpText = `
Welcome to the Pokedex!
Usage:

`

func NewCommands() map[string]*CliCommand {
	commands := make(map[string]*CliCommand)
	commands["help"] = newHelpCommand()
	commands["exit"] = newExitCommand()
	for k, v := range commands {
		helpText += fmt.Sprintf("%s: %s\n", k, v.description)
	}
	return commands
}
