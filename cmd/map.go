package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Flarenzy/Pokedex/internal/config"
)

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationArea struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

func getFromAPI(url string, c *config.Config) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.Logger.Error("Error creating request: ", "error", err)
		return []byte{}, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Error("Error making request: ", "error", err)
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			c.Logger.Error("Error closing body: ", "error", err)
			panic("error closing body")
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Error("Error reading response: ", "url", url, "error", err)
		return []byte{}, err
	}
	return body, nil
}

func getLocationArea(c *config.Config, url string) error {
	cachedBody, err := c.Cache.Get(url)
	var body []byte
	if err != nil {
		body, err = getFromAPI(url, c)
		if err != nil {
			return err
		}
		c.Logger.Debug("Adding key to cache: ", "url", url)
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

	var locationsArea LocationArea
	err = json.Unmarshal(body, &locationsArea)
	if err != nil {
		c.Logger.Error("Error parsing response: ", "error", err)
		return err
	}

	for _, location := range locationsArea.Results {
		_, err = fmt.Fprintln(c.Out, location.Name)
		if err != nil {
			c.Logger.Error("Error writing response: ", "url", url, "error", err)
			return err
		}
	}

	c.Next = locationsArea.Next
	c.Previous = locationsArea.Previous
	return nil
}

func commandMap(c *config.Config) error {
	url := c.Next
	err := getLocationArea(c, url)
	if err != nil {
		return err
	}
	return nil
}

func commandMapb(c *config.Config) error {
	url := c.Previous
	if url == "" {
		_, err := fmt.Fprintln(c.Out, "you're on the first page")
		if err != nil {
			c.Logger.Error("Error getting response: ", "error", err)
			return err
		}
		return nil
	}
	err := getLocationArea(c, url)
	if err != nil {
		c.Logger.Error("Error getting location area: ", "error", err)
		return err
	}
	return nil
}

func newMapCommand() *CliCommand {
	return &CliCommand{
		name:        "map",
		description: "Display the next location of a map",
		Callback:    commandMap,
	}
}

func newMapbCommand() *CliCommand {
	return &CliCommand{
		name:        "mapb",
		description: "Display the previous location of a map",
		Callback:    commandMapb,
	}
}
