package cli

import "fmt"

func commandHelp(cfg *config) error {
	fmt.Print("Welcome to the Pokedex! \nUsage:\n")
	for _, command := range cfg.command {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
