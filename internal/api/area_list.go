package api

import (
	"fmt"
	"io"
	"encoding/json"
)

type AreasList struct {
	Count    int    
	Next     *string 
	Previous *string    
	Results  []struct {
		Name string
		URL  string 
	}
}

func (c *Client) GetAreas(nextPage *string) (AreasList, error) {
	url := c.baseURL + "/location-area/"
	if nextPage != nil {
		url = *nextPage
	}
	
	body, ok := c.Cache.Get(url)
	if !ok {
		// send the Response
		res, err := c.httpClient.Get(url)
		if err != nil {
			return AreasList{}, err
		}
	
		defer res.Body.Close()
		// get response Body
		data, err := io.ReadAll(res.Body)
	
		if err != nil {
			return AreasList{}, err
		}
		
		if res.StatusCode > 299 {
			return AreasList{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, data)
		}
		fmt.Println("Cache miss: making network request...")
		body = data
	}else {
		fmt.Println("Cache hit!")
	}
	
	//parse res.Body
	var areas AreasList
	if err := json.Unmarshal(body, &areas); err != nil {
		return AreasList{}, err
	}
	if !ok {
		c.Cache.Add(url, body)
	}
	return areas, nil
}
