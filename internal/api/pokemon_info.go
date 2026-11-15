package api

import (
	"fmt"
	"io"
	"encoding/json"
)

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (c *Client) GetPokemonInfo(name string) (Pokemon, error) {
	url := c.baseURL + "/pokemon/" + name

	body, ok := c.Cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}
		if res.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("Error. Status code: %d and Body: %s\n", res.StatusCode, data)
		}
		body = data
		fmt.Println("Cache miss: sending network query")
	}else {		
		fmt.Println("Cache hit!")
	}
	
	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, err
	}
	if !ok {
		c.Cache.Add(url, body)
	}
	return pokemon, nil
}
