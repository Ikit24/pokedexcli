package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"math/rand"
)

func commandCatch(cfg *config, pokemon_name []string) error {
	if len(pokemon_name) == 0 {
		return fmt.Errorf("Must provide pokemon name in order to catch")
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

	var baseXP_low = 100
	var baseXP_medium = 150
	var baseXP_high = 300

	var catch_Chance int

	if apiResponse.BaseExperience <= baseXP_low {
		catch_Chance = 80
	} else if apiResponse.BaseExperience <= baseXP_medium {
		catch_Chance = 50
	} else if apiResponse.BaseExperience <= baseXP_high {
		catch_Chance = 20
	} else {
		catch_Chance = 5
	}

	randomRoll := rand.Intn(100) + 1
	if randomRoll <= catch_Chance {
		fmt.Printf("%s was caught!\n", pokemon_name[0])
		cfg.Caught[apiResponse.Name] = apiResponse
	} else {
		fmt.Printf("%s escaped!\n", pokemon_name[0])
	}
	return nil
}
