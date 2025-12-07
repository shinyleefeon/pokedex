package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
)

func cleanInput(text string) []string {
	var result []string
	for _, arg := range strings.Fields(text) {
		// Add cleaning logic here, e.g., trimming spaces, converting to lowercase, etc.
		result = append(result, strings.ToLower(strings.TrimSpace(arg)))
	}
	return result
}











func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		args := cleanInput(input)
		fmt.Printf("Your command was: %v\n", args[0])
		// Process args here
	}
}

