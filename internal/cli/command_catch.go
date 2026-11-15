package cli

import (
	"fmt"
	"math/rand"
)

func commandCatch(cfg *config, args ...string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}else {
		return fmt.Errorf("You must specify pokemon to catch")
	}
	pokemon, err := cfg.apiClient.GetPokemonInfo(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	chanceCatch := rand.Intn(pokemon.BaseExperience)

	if chanceCatch > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}
	cfg.caughtPokemons[name] = pokemon
	fmt.Printf("%s was caught!\n", pokemon.Name)
	
	return nil
}
