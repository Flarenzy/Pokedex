package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func getFromAPI(url string, config *Config) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		config.Logger.Error("Error creating request: ", "error", err)
		return []byte{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		config.Logger.Error("Error making request: ", "error", err)
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			config.Logger.Error("Error closing body: ", "error", err)
			panic("error closing body")
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("Error reading response: ", "url", url, "error", err)
		return []byte{}, err
	}
	return body, nil
}

func getLocationArea(config *Config, url string) error {
	cachedBody, err := config.cache.Get(url)
	var body []byte
	if err != nil {
		body, err = getFromAPI(url, config)
		if err != nil {
			return err
		}
	} else {
		body = cachedBody
		//fmt.Println("cache hit")
		config.Logger.Debug("Cache hit", "url", url)
	}

	var locationsArea LocationArea
	err = json.Unmarshal(body, &locationsArea)
	if err != nil {
		config.Logger.Error("Error parsing response: ", "error", err)
		return err
	}

	for _, location := range locationsArea.Results {
		fmt.Println(location.Name)
	}
	config.Logger.Debug("Adding key to cache: ", "url", url)
	err1 := config.cache.Add(url, body)
	if err1 != nil {
		if err1.Error() != fmt.Sprintf("key already exists: %v", url) {
			return err1
		}
	}
	config.Next = locationsArea.Next
	config.Previous = locationsArea.Previous
	return nil
}

func commandMap(config *Config) error {
	url := config.Next
	err := getLocationArea(config, url)
	if err != nil {
		return err
	}
	return nil
}

func commandMapb(config *Config) error {
	url := config.Previous
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	err := getLocationArea(config, url)
	if err != nil {
		config.Logger.Error("Error getting location area: ", "error", err)
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
