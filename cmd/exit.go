package cmd

import (
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/config"
	"os"
)

func commandExit(c *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	c.Cache.Done()
	os.Exit(0)
	return nil
}

func newExitCommand() *CliCommand {
	return &CliCommand{
		name:        "exit",
		description: "Exit Pokedex",
		Callback:    commandExit,
	}
}
