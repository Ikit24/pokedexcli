package commands

import(
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandExplore(cfg *config.Config, location_area_name []string) error {
	if cfg.Next == "" || cfg.Next == "https://pokeapi.co/api/v2/location-area/" {
		return fmt.Errorf("You must map the area first! Use the 'map' command.")
	}

	if len(location_area_name) == 0 {
		return fmt.Errorf("Must provide location area name!")
	}
	var explore_URL = "https://pokeapi.co/api/v2/location-area/" + location_area_name[0] + "/"
	fmt.Println("Exploring...")

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

	for _, encounter := range apiResponse.PokemonEncounters {
		pokemonName := encounter.Pokemon.Name
		specificPokemonURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"

		cachedPokemonData, ok := cfg.Cache.Get(specificPokemonURL)
		if !ok {
			res, httpErr := http.Get(specificPokemonURL)
			if httpErr != nil {
				return httpErr
			}
			defer res.Body.Close()
			body, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			cfg.Cache.Add(specificPokemonURL, body)
			if res.StatusCode > 299 {
				return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			}
		} else {
			body = cachedPokemonData
			err = nil
		}
		var pokemonResponse pokeapi.BattlePokemon
		err = json.Unmarshal(body, &pokemonResponse)
		if err != nil {
			return fmt.Errorf("JSON unmarshal failed: %w", err)
		}
		cfg.Battle[pokemonName] = pokemonResponse
	}
	fmt.Println("Exploring", location_area_name[0], "...")
	fmt.Println()
	fmt.Println("Found Pokemon:")

	for _, encounter := range apiResponse.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
