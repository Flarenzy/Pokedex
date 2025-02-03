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
	fmt.Println("Your Pokedex:")
	for _, pokemon := range allPokemon {
		fmt.Printf("  - %v\n", pokemon.Name)
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
