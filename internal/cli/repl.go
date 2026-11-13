package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/danl1k3/pokedex/internal/api"
	"time"
)

func Start() {
	client := api.NewClient(5 * time.Second)

	cfg := config{
		apiClient: client,
		command: getCommands(),
	}
	startRepl(&cfg)
}

func startRepl(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := cfg.command[commandName]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}



