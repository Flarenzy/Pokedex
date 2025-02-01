package cmd

import "fmt"

func commandHelp(config *Config) error {
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
