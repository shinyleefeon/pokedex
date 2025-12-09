package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"github.com/shinyleefeon/pokedex/internal/pokecache"
	"time"
	"math/rand"
)

type LocationAreaResponse struct {
	Results []LocationArea `json:"results"`
	Count  int            `json:"count"`
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

type LocationArea struct {
	Name string `json:"name"`
	Location struct {
		Name     string
		Location string
	}
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []Stat `json:"stats"`
	Types []Type `json:"types"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type Type struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"type"`
}

// contains checks if a string is in a slice
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

var Offset int = 0
var cache = pokecache.NewCache(make(map[string]pokecache.CacheEntry), 10*time.Minute)

func MapCommand(args []string) error {
	req := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=20&offset=%d", Offset)
	var body []byte
	data, found := cache.Get(req)
	if found {
		body = data
		//println("USED CACHE")

	} else {
		resp, err := http.Get(req)
		if err != nil {
			return fmt.Errorf("failed to fetch locations: %v", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cache.Add(req, body)
	}
	var result = LocationAreaResponse{}	
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, location := range result.Results {
		fmt.Printf("Location: %s\n", location.Name)
	}

	Offset += 20

	return nil
}

func MapbackCommand(args []string) error {
	if Offset > 20 {
		Offset -= 40
		return MapCommand(nil)
	}
	fmt.Println("you're on the first page")
	return nil
}

func ExploreCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("please provide a location area name")
	}
	locationAreaName := args[1]
	req := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationAreaName)

	var body []byte
	data, found := cache.Get(req)
	if found {
		body = data
		//println("USED CACHE")
	} else {
		resp, err := http.Get(req)
		if err != nil {
			return fmt.Errorf("failed to fetch location area: %v", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cache.Add(req, body)
	}

	// Further processing of the location area data can be done here
	var pokemonList = make([]string, 0)
	var result = LocationAreaResponse{}	
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, encounter := range result.Encounters {
		pokemonList = append(pokemonList, encounter.Pokemon.Name)
	}
	fmt.Printf("Exploring %v...,\n", locationAreaName)
	fmt.Printf("Found Pokémon: \n")
	for _, pokemon := range pokemonList {
		fmt.Printf("- %s\n", pokemon)
	}

	return nil
}

func CaculateChance(difficulty int) float64 {
	if difficulty < 1 { difficulty = 1}
	if difficulty > 700 { difficulty = 700}
	const (
		minDifficulty = 1
		maxDifficulty = 700
		minChance     = 0.01
		maxChance     = 0.90
	)

	// Linear interpolation to calculate chance based on difficulty
	slope := (minChance - maxChance) / (maxDifficulty - minDifficulty)
	probability := maxChance + (float64(difficulty)-minDifficulty)*slope
	//println("Calculated catch probability:", probability)
	return probability
}

func RollCatch(chance float64) bool {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Float64() < chance
}


func Catch(args []string) (bool, error) {
	if len(args) < 2 {
		return false, fmt.Errorf("please provide a pokemon name to catch")
	}
	pokemonName := args[1]
	req := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	var body []byte
	data, found := cache.Get(req)
	if found {
		body = data
		//println("USED CACHE")
	} else {
		resp, err := http.Get(req)
		if err != nil {
			return false, fmt.Errorf("failed to fetch pokemon: %v", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("failed to read response body: %v", err)
		}
		cache.Add(req, body)
	}

	var result Pokemon
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse JSON: %v", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", result.Name)
	chance := CaculateChance(result.BaseExperience)
	if RollCatch(chance) {
		fmt.Printf("%s was caught!\n", result.Name)
		return true, nil
	} else {
		fmt.Printf("%s escaped!\n", result.Name)
	}
	return false, nil
	
	

	
}

func Inspect(pokedex []string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("please provide a pokemon name to inspect")
	}
	pokemonName := args[1]
	if !contains(pokedex, pokemonName) {
		return fmt.Errorf("you don't have %s in your pokedex", pokemonName)
	}
	req := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	var body []byte
	data, found := cache.Get(req)
	if found {
		body = data
		//println("USED CACHE")
	} else {
		resp, err := http.Get(req)
		if err != nil {
			return fmt.Errorf("failed to fetch pokemon: %v", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cache.Add(req, body)
	}

	var result Pokemon
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}
	fmt.Printf("Name: %s\n", result.Name)
	fmt.Printf("Height: %d\n", result.Height)
	fmt.Printf("Weight: %d\n", result.Weight)
	fmt.Printf("Stats:\n")
	for _, s := range result.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range result.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}	
		return nil
}


