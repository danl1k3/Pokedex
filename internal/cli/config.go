package cli

import "github.com/danl1k3/pokedex/internal/api"

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
	}
}

type config struct {
	apiClient 	api.Client
	nextArea 	*string
	previousArea 	*string
	command 	map[string]cliCommand
}


type cliCommand struct {
	name        string
	description string
	callback    func(*config) error

}

