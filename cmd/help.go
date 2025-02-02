package cmd

import (
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/config"
)

func commandHelp(c *config.Config) error {
	fmt.Print(helpText)
	return nil
}

func newHelpCommand() *CliCommand {
	return &CliCommand{
		name:        "help",
		description: "Displays a help message",
		Callback:    commandHelp,
	}
}
