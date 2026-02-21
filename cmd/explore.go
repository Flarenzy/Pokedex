package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/config"
)

type PokemonInLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string            `json:"name"`
		Language map[string]string `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel          int           `json:"min_level"`
				MaxLevel          int           `json:"max_level"`
				ConditionalValues []interface{} `json:"conditional_values"`
				Chance            int           `json:"chance"`
				Method            struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
			}
		}
	} `json:"pokemon_encounters"`
}

func getPokemonInArea(c *config.Config, url string) error {
	body, err := getBodyWithCache(c, url)
	if err != nil {
		return err
	}
	var pokemonInLocation PokemonInLocation
	err = json.Unmarshal(body, &pokemonInLocation)
	if err != nil {
		c.Logger.Error("Error parsing response: ", "error", err)
		return err
	}
	for i, pokemon := range pokemonInLocation.PokemonEncounters {
		_, err = fmt.Fprintf(c.Out, "Pokemon #%v: %v\n", i+1, pokemon.Pokemon.Name)
		if err != nil {
			c.Logger.Error("Error writing response: ", "error", err)
			return err
		}
	}
	c.Logger.Debug("Adding key to cache: ", "url", url)
	return nil
}

func commandExplore(c *config.Config) error {
	if len(c.Args) == 0 {
		c.Logger.Info("No command to explore")
		return ErrNoAreaToExplore
	}
	for _, arg := range c.Args {
		baseURL := c.AreaURL
		if baseURL == "" {
			baseURL = internal.FirstURL
		}
		url := baseURL + arg
		_, err := fmt.Fprintln(c.Out, "Exploring area: ", arg)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(c.Out, "===================================")
		if err != nil {
			return err
		}
		err = getPokemonInArea(c, url)
		if err != nil {
			c.Logger.Error("Error getting pokemon in area: ", "error", err)
		}
		_, err = fmt.Fprintln(c.Out, "===================================")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(c.Out)
		if err != nil {
			return err
		}
	}
	return nil
}

func newExploreCommand() *CliCommand {
	return &CliCommand{
		name:        "explore",
		description: "Explore an pokemon in an certain area, area are provided by the map and mapb commands.",
		Callback:    commandExplore,
	}
}
