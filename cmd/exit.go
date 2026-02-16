package cmd

import (
	"fmt"

	"github.com/Flarenzy/Pokedex/internal/config"
)

func commandExit(c *config.Config) error {
	_, err := fmt.Fprintln(c.Out, "Closing the Pokedex... Goodbye!")
	if err != nil {
		c.Logger.Error("unable to close the Pokedex", "error", err)
		return err
	}
	return ErrStop
}

func newExitCommand() *CliCommand {
	return &CliCommand{
		name:        "exit",
		description: "Exit Pokedex",
		Callback:    commandExit,
	}
}
