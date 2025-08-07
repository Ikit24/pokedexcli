package main

import(
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)

func commandExplore(cfg *config, location_area_name []string) error {
	if len(location_area_name) == 0 {
		return fmt.Errorf("Must provide loaction area name")
	}
	var explore_URL = "https://pokeapi.co/api/v2/location-area/" + location_area_name[0] + "/"

	var body []byte
	var err error

	cachedData, ok := cfg.Cache.Get(explore_URL)
	if !ok {
		res, httpErr := http.Get(explore_URL)
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(explore_URL, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	} else {
		body = cachedData
		err = nil
	}
	var apiResponse pokeapi.LocationAreaDetails

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	fmt.Println("Exploring", location_area_name[0], "...")
	fmt.Println("Found Pokemon:")

	for _, encounter := range apiResponse.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
