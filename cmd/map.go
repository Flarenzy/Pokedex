package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

func getLocationArea(url string) (LocationArea, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request: ", err)
		return LocationArea{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error making request: ", err)
		return LocationArea{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing body: ", err)
			panic("error closing body")
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response: ", err)
		return LocationArea{}, err
	}
	var locationsArea LocationArea
	err = json.Unmarshal(body, &locationsArea)
	if err != nil {
		slog.Error("Error parsing response: ", err)
		return LocationArea{}, err
	}

	for _, location := range locationsArea.Results {
		fmt.Println(location.Name)
	}
	return locationsArea, nil
}

func commandMap(config *Config) error {
	url := config.Next
	la, err := getLocationArea(url)
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
	la, err := getLocationArea(url)
	if err != nil {
		slog.Error("Error getting location area: ", err)
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
