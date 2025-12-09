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
	callback    func(args []string) error
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
		"explore": {
			name:        "explore",
			description: "returns a list of pokemon at the specified location",
			callback:    pokeapi.ExploreCommand,
		},
		"catch": {
			name:        "catch",
			description: "Catch the specified pokemon",
			callback:    CatchCommand,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect your caught pokemon",
			callback:    InspectCommand,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all caught pokemon in your pokedex",
			callback:    PokedexPrint,
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

func commandExit(args []string) error {
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

func commandHelp(args []string) error {
	return showCommands(commands)
}

var Pokedex = make([]string, 0)

func CatchCommand(args []string) error {
	caught, err := pokeapi.Catch(args)
	if err != nil {
		return err
	}
	if caught {
		Pokedex = append(Pokedex, args[1])
	}
	return nil
}

func  InspectCommand(args []string) error {
	if len(Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokémon first!")
		return nil
	}
	err := pokeapi.Inspect(Pokedex, args)
	return err
}

func PokedexPrint(args []string) error {
	if len(Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokémon first!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range Pokedex {
		fmt.Printf("- %s\n", pokemon)
	}
	return nil
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
			commandExit(nil)
		case "help":
			commandHelp(nil)
		case "map":
			err := pokeapi.MapCommand(nil)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "mapb":
			err := pokeapi.MapbackCommand(nil)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "explore":
			err := pokeapi.ExploreCommand(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "catch":
			err := CatchCommand(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "inspect":
			err := InspectCommand(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "pokedex":
			err := PokedexPrint(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println("Unknown command. Please try again.")
		}
	}
}

