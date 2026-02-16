package cmd

import (
	"fmt"

	"github.com/Flarenzy/Pokedex/internal/config"
)

func commandPokedex(c *config.Config) error {
	allPokemon := c.Pokedex.GetAllPokemon()
	if len(allPokemon) == 0 {
		return fmt.Errorf("no Pokedex found")
	}
	_, err := fmt.Fprintln(c.Out, "Your Pokedex:")
	if err != nil {
		c.Logger.Error("Failed to write to output", "error", err)
		return err
	}
	for _, pokemon := range allPokemon {
		_, err = fmt.Fprintf(c.Out, "  - %v\n", pokemon.Name)
		if err != nil {
			c.Logger.Error("Failed to write to output", "error", err)
			return err
		}
	}
	return nil
}

func newPokedexCommand() *CliCommand {
	return &CliCommand{
		name:        "pokedex",
		description: `Displays all pokemon in pokedex`,
		Callback:    commandPokedex,
	}
}
