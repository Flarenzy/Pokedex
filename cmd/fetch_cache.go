package cmd

import (
	"errors"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
)

func getBodyWithCache(c *config.Config, url string) ([]byte, error) {
	cachedBody, err := c.Cache.Get(url)
	if err == nil {
		c.Logger.Debug("Cache hit", "url", url)
		return cachedBody, nil
	}

	body, err := getFromAPI(url, c)
	if err != nil {
		return nil, err
	}

	c.Logger.Debug("Adding key to cache: ", "url", url)
	err = c.Cache.Add(url, body)
	if err != nil && !errors.Is(err, pokecache.ErrKeyExists) {
		return nil, err
	}

	return body, nil
}
