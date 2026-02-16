package http

import "net/http"

type HTTPClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

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
