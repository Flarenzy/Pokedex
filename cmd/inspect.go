package cmd

import (
	"errors"
	"fmt"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
)

func commandInspect(c *config.Config) error {
	if len(c.Args) == 0 {
		c.Logger.Info("No pokemon to inspect")
		return ErrNoPokemonToInspect
	}
	for _, arg := range c.Args {
		p, err := c.Pokedex.GetPokemonByName(arg)
		if err != nil {
			if errors.Is(err, pokedex.ErrPokemonNotFound) {
				_, err = fmt.Fprintf(c.Out, "you have not caught that pokemon")
				if err != nil {
					c.Logger.Error("Error writing response: ", "error", err)
					return err
				}
				c.Logger.Error("You have not caught that pokemon", "err", err)
				return nil
			}
			c.Logger.Error("Error writing response: ", "error", err)
			return err
		}
		c.Logger.Debug("Printing info about pokemon: ", "name", p.Name)
		_, err = fmt.Fprintf(c.Out, "Name: %v\nHeight: %v\nWeight: %v\n", p.Name, p.Height, p.Weight)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(c.Out, "Stats:")
		if err != nil {
			return err
		}
		for _, stat := range p.Stats {
			_, err = fmt.Fprintf(c.Out, "  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
			if err != nil {
				c.Logger.Error("Error writing response: ", "error", err)
				return err
			}
		}
		_, err = fmt.Fprintln(c.Out, "Types:")
		if err != nil {
			return err
		}
		for _, t := range p.Types {
			_, err = fmt.Fprintf(c.Out, "  -%v\n", t.Type.Name)
			if err != nil {
				c.Logger.Error("Error writing response: ", "error", err)
				return err
			}
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
