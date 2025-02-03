package cmd

import (
	"errors"
	"fmt"
	"github.com/Flarenzy/Pokedex/internal/config"
)

func commandInspect(c *config.Config) error {
	if len(c.Args) == 0 {
		c.Logger.Info("No pokemon to inspect")
		return errors.New("no pokemon to inspect")
	}
	for _, arg := range c.Args {
		p, err := c.Pokedex.GetPokemonByName(arg)
		if err != nil {
			if err.Error() == "pokemon not found" {
				fmt.Printf("you have not caught that pokemon")
				c.Logger.Error("You have not caught that pokemon", "err", err)
				return nil
			} else {
				return err
			}
		}
		c.Logger.Debug("Printing info about pokemon: ", "name", p.Name)
		fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\n", p.Name, p.Height, p.Weight)
		fmt.Println("Stats:")
		for _, stat := range p.Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range p.Types {
			fmt.Printf("  -%v\n", t.Type.Name)
		}
	}
	return nil
}

func newInspectCommand() *CliCommand {
	return &CliCommand{
		name:        "inspect",
		description: `Inspect pokemon.`,
		Callback:    commandInspect,
	}
}
