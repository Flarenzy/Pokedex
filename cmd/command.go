package cmd

import "fmt"

type CliCommand struct {
	name        string
	description string
	Callback    func(config *Config) error
}

type Config struct {
	Next     string
	Previous string
}

func NewConfig() *Config {
	return &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
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
	commands["map"] = NewMapCommand()
	commands["mapb"] = NewMapbCommand()
	for k, v := range commands {
		helpText += fmt.Sprintf("%s: %s\n", k, v.description)
	}
	return commands
}
