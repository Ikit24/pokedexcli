package commands

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandMap(cfg *config.Config, _ []string) error {
	fmt.Println("Fetching locations...")
	var body []byte
	var err error

	cachedData, ok := cfg.Cache.Get(cfg.Next)
	if !ok {
		res, httpErr := http.Get(cfg.Next)
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(cfg.Next, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	} else {
		body = cachedData
		err = nil
	}

	var apiResponse pokeapi.LocationAreasResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	if apiResponse.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *apiResponse.Previous
	}

	if apiResponse.Next == "" {
		cfg.Next = ""
	} else {
		cfg.Next = apiResponse.Next
	}

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}
	return nil
	}

func CommandMapb(cfg *config.Config, _ []string) error {
	if cfg.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	var body []byte
	var err error

	cachedData, ok := cfg.Cache.Get(cfg.Previous)
	if !ok {
		res, httpErr := http.Get(cfg.Previous)
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(cfg.Previous, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	} else {
		body = cachedData
		err = nil
	}

	var apiResponse pokeapi.LocationAreasResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	if apiResponse.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *apiResponse.Previous
	}

	if apiResponse.Next == "" {
		cfg.Next = ""
	} else {
		cfg.Next = apiResponse.Next
	}

	for _, result := range apiResponse.Results {
		fmt.Println(result.Name)
	}
	return nil
}
