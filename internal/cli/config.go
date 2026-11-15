package cli

import (
	"github.com/danl1k3/pokedex/internal/api"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name: 	     "map",
			description: "Lists next areas of Pokemon World",
			callback:    commandMap,
		},
		"mapb": {
			name:	     "mapb",
			description: "Lists previous areas of Pokemon World",
			callback:    commandMapb,
		},
		"explore": {
			name: 	     "explore",
			description: "Lists all pokemons name in specified area",
			callback:    commandExplore,
		},
		"catch": {
			name:	     "catch",
			description: "Give abillity to catch the pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:	     "inspect",
			description: "Gives info about caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:	     "pokedex",
			description: "Lists all caught pokemons",
			callback:    commandPokedex,
		},
	}
}

type config struct {
	apiClient 	api.Client
	nextArea 	*string
	previousArea 	*string
	caughtPokemons  map[string]api.Pokemon
	command 	map[string]cliCommand
}


type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error

}

