package cmd

import (
	"fmt"

	"github.com/Flarenzy/Pokedex/internal/config"
)

func commandHelp(c *config.Config) error {
	_, err := fmt.Fprintln(c.Out, helpText)
	if err != nil {
		c.Logger.Error("unable to write help message", "error", err.Error())
		return err
	}
	return nil
}

func newHelpCommand() *CliCommand {
	return &CliCommand{
		name:        "help",
		description: "Displays a help message",
		Callback:    commandHelp,
	}
}
