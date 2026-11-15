package cli 

import (
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}else{
		fmt.Errorf("You must specify pokemon")
	}

	pokemon, exists := cfg.caughtPokemons[name]
	if !exists {
		return fmt.Errorf("You havn't caught that pokemon yet")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}	

	return nil
}
