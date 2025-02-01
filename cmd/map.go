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
		config.Logger.Error("Error creating request: ", err)
		return []byte{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		config.Logger.Error("Error making request: ", err)
		return []byte{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			config.Logger.Error("Error closing body: ", err)
			panic("error closing body")
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Error("Error reading response: ", err)
		return []byte{}, err
	}
	return body, nil
}

func getLocationArea(config *Config, url string) (LocationArea, error) {
	cachedBody, err := config.cache.Get(url)
	var body []byte
	if err != nil {
		body, err = getFromAPI(url, config)
		if err != nil {
			return LocationArea{}, err
		}
	} else {
		body = cachedBody
	}

	var locationsArea LocationArea
	err = json.Unmarshal(body, &locationsArea)
	if err != nil {
		config.Logger.Error("Error parsing response: ", err)
		return LocationArea{}, err
	}

	for _, location := range locationsArea.Results {
		fmt.Println(location.Name)
	}
	config.Logger.Info("Adding key to cache: ", url)
	err1 := config.cache.Add(url, body)
	if err1 != nil {
		return LocationArea{}, err1
	}
	return locationsArea, nil
}

func commandMap(config *Config) error {
	url := config.Next
	la, err := getLocationArea(config, url)
	if err != nil {
		return err
	}
	config.Next = la.Next
	config.Previous = la.Previous
	return nil
}

func commandMapb(config *Config) error {
	url := config.Previous
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	la, err := getLocationArea(config, url)
	if err != nil {
		config.Logger.Error("Error getting location area: ", err)
		return err
	}
	config.Next = la.Next
	config.Previous = la.Previous
	return nil
}

func NewMapCommand() *CliCommand {
	return &CliCommand{
		name:        "map",
		description: "Display the next location of a map",
		Callback:    commandMap,
	}
}

func NewMapbCommand() *CliCommand {
	return &CliCommand{
		name:        "mapb",
		description: "Display the previous location of a map",
		Callback:    commandMapb,
	}
}
