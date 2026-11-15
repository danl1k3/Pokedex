package api

import (
	"net/http"
	"time"
	"github.com/danl1k3/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	baseURL	   string
	Cache	   *pokecache.Cache
}

func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		baseURL: "https://pokeapi.co/api/v2",
		Cache: pokecache.NewCache(cacheInterval),
	}
}

func NewClientWithURL(timeout time.Duration, url string, cacheInterval time.Duration) Client{
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		baseURL: url,
		Cache: pokecache.NewCache(cacheInterval),
	}
}
