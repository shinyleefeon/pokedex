package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
	"github.com/shinyleefeon/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show 20 locations from the Pokedex",
			callback:    pokeapi.MapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous 20 locations from the Pokedex",
			callback:    pokeapi.MapbackCommand,
		},
	}
}

func cleanInput(text string) []string {
	var result []string
	for _, arg := range strings.Fields(text) {
		// Add cleaning logic here, e.g., trimming spaces, converting to lowercase, etc.
		result = append(result, strings.ToLower(strings.TrimSpace(arg)))
	}
	return result
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func showCommands(registry map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range registry {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandHelp() error {
	return showCommands(commands)
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		args := cleanInput(input)
		fmt.Printf("Your command was: %v\n", args[0])
		switch args[0] {
		case "exit":
			commandExit()
		case "help":
			commandHelp()
		case "map":
			err := pokeapi.MapCommand()
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "mapb":
			err := pokeapi.MapbackCommand()
			if err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println("Unknown command. Please try again.")
		}
	}
}

