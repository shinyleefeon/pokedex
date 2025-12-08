package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
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

func MapCommand() error {
	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=20&offset=%d", Offset))
	if err != nil {
		return fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
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