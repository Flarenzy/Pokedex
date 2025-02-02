package cmd

import (
	"encoding/json"
	"errors"
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
	cachedBody, err := c.Cache.Get(url)
	var body []byte
	if err != nil {
		body, err = getFromAPI(url, c)
		if err != nil {
			return err
		}
		err1 := c.Cache.Add(url, body)
		if err1 != nil {
			if err1.Error() != fmt.Sprintf("key already exists: %v", url) {
				return err1
			}
		}
	} else {
		body = cachedBody
		c.Logger.Debug("Cache hit", "url", url)
	}
	var pokemonInLocation PokemonInLocation
	err = json.Unmarshal(body, &pokemonInLocation)
	if err != nil {
		c.Logger.Error("Error parsing response: ", "error", err)
		return err
	}
	for i, pokemon := range pokemonInLocation.PokemonEncounters {
		fmt.Printf("Pokemon #%v: %v\n", i+1, pokemon.Pokemon.Name)
	}
	c.Logger.Debug("Adding key to cache: ", "url", url)
	return nil
}

func commandExplore(c *config.Config) error {
	if len(c.Args) == 0 {
		c.Logger.Info("No command to explore")
		return errors.New("no command to explore")
	}
	//wg := sync.WaitGroup{}
	//wg.Add(len(c.Args))
	for _, arg := range c.Args {
		url := internal.FirstURL + arg // TODO kako uzeti pokemon area
		fmt.Println("Exploring area: ", arg)
		fmt.Println("===================================")
		err := getPokemonInArea(c, url)
		if err != nil {
			c.Logger.Error("Error getting pokemon in area: ", "error", err)
		}
		fmt.Println("===================================")
		fmt.Println()
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
