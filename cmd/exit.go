package cmd

import (
	"fmt"
	"os"
)

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
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
