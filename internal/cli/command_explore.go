package cli

import "fmt"

func commandExplore(cfg *config, args  ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location name")
	}
	name := args[0]
	fmt.Printf("Exploring %s\n", name)

	response, err := cfg.apiClient.GetAreaInfo(name)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	
	fmt.Println("Found Pokemons:")
	for _, encounter := range response.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}
