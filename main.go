package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"github.com/Ikit24/pokedexcli/internal/config"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/pokecache"
)

func main() {
	cfg, err := loadOrCreateConfig()
	if err != nil {
		fmt.Printf("Error loading save: %v\n", err)
		return
	}
	startRepl(cfg)
}

func loadOrCreateConfig() (config.Config, error) {
	var cfg config.Config

	save, err := ioutil.ReadFile("pokedex.json")
	if err != nil {
		cfg.Caught = make(map[string]pokeapi.BattlePokemon)
		cfg.ExploredAreas = []string{}
	} else {
		var saveData struct {
			Caught        map[string]pokeapi.BattlePokemon `json:"caught"`
			ExploredAreas []string                        `json:"explored_areas"`
		}
		err := json.Unmarshal(save, &saveData)
		if err != nil {
			fmt.Println("Save file corrupted, creating a new save...")
			os.Remove("pokedex.json")
			cfg.Caught = make(map[string]pokeapi.BattlePokemon)
			cfg.ExploredAreas = []string{}
		} else {
			cfg.Caught = saveData.Caught
			cfg.ExploredAreas = saveData.ExploredAreas
		}
	}
	cfg.Next = "https://pokeapi.co/api/v2/location-area/"
	cfg.Previous = ""
	cfg.Cache = pokecache.NewCache(1 * time.Minute)
	cfg.MyMap = getCommands(&cfg)
	cfg.Battle = make(map[string]pokeapi.BattlePokemon)

	return cfg, nil
}
