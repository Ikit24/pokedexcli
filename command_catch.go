package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)

func commandCatch(cfg *config, pokemon_name []string) error {
	if len(pokemon_name) == 0 {
		return fmt.ErrorF("Must provide pokemon name in order to catch")
	}
	
	var catch_URL = "https://pokeapi.co/api/v2/pokemon/" + pokemon_name[0] + "/"

	var body []byte
	var err error

	cachedData, ok := cfg.Cache.Get(catch_URL)
	if !ok {
		res, httpErr := http.Get(catch_URL)
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(catch_URL, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	} else {
		body = cachedData
		err = nil
	}
	var apiResponse pokeapi.Pokemon

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...", pokemon_name[0])
	fmt.Println()
}
