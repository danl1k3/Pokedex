package cli 

import "fmt"


func commandMapb(cfg *config) error {
	if cfg.previousArea == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	response, err := cfg.apiClient.GetAreas(cfg.previousArea)
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
