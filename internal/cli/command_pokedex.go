package cli

import (
	"fmt"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("Pokedex takes no args!")
	}
	if len(cfg.caughtPokemons) == 0 {
		fmt.Println("Your Pokedex is empty")
	}else {
		fmt.Println("Your Pokedex:")
		for _, pokemon := range cfg.caughtPokemons {
			fmt.Printf(" - %s\n", pokemon.Name)
		}
	}
	return nil
}
