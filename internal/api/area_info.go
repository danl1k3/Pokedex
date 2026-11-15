package api

import (
	"fmt"
	"io"
	"encoding/json"
)

type AreaInfo struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
				Chance   int `json:"chance"`
				Method   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				ConditionValues []any `json:"condition_values"`
				PokemonDetails  any   `json:"pokemon_details"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetAreaInfo(name string) (AreaInfo, error) {
	url := c.baseURL + "/location-area/" + name
	
	body, ok := c.Cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return AreaInfo{}, err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return AreaInfo{}, err
		}
		if res.StatusCode > 299 {
			return AreaInfo{}, fmt.Errorf("Error. Status code: %d and Body: %s\n", res.StatusCode, data)
		}
		body = data
		fmt.Println("Cache miss: sending network query")
	}else {
		fmt.Println("Cache hit!")
	}
	
	var areaInfo AreaInfo
	if err := json.Unmarshal(body, &areaInfo); err != nil {
		return AreaInfo{}, err
	}
	if !ok {
		c.Cache.Add(url, body)
	}
	return areaInfo, nil
}
