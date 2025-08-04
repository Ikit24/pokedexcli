package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)

func commandMap(cfg *config) error {
	res, err := http.Get(cfg.Next)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var apiResponse pokeapi.LocationAreasResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarhsal failed: %w", err)
	}

	if apiResponse.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *apiResponse.Previous
	}

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}
	return nil
	}

func commandMapb(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(cfg.Previous)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var apiResponse pokeapi.LocationAreasResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarhsal failed: %w", err)
	}

	if apiResponse.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *apiResponse.Previous
	}

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}
	return nil
}
