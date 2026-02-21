package http

import "net/http"

import "github.com/Flarenzy/Pokedex/internal/domain"

type HTTPClienter = domain.HTTPClient

type DefaultHTTPClient struct {
	c *http.Client
}

func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		c: &http.Client{},
	}
}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}
