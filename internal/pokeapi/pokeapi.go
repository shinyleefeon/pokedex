package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"github.com/shinyleefeon/pokedex/internal/pokecache"
	"time"
)

type LocationAreaResponse struct {
	Results []LocationArea `json:"results"`
	Count  int            `json:"count"`
}

type LocationArea struct {
	Name string `json:"name"`
	Location struct {
		Name     string
		Location string
	}
}

var Offset int = 0
var cache = pokecache.NewCache(make(map[string]pokecache.CacheEntry), 10*time.Minute)

func MapCommand() error {
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

func MapbackCommand() error {
	if Offset > 20 {
		Offset -= 40
		return MapCommand()
	}
	fmt.Println("you're on the first page")
	return nil
}