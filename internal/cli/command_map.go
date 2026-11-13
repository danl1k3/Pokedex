package cli

import (
	"fmt"
)

func commandMap(cfg *config) error {
	response, err := cfg.apiClient.GetAreas(cfg.nextArea)
	if err != nil {
		return err
	}

	cfg.nextArea = response.Next
	cfg.previousArea = response.Previous

	for _, res := range response.Results {
		fmt.Println(res.Name)
	}

	return nil
}
