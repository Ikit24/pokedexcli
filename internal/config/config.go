package config

import (
    "github.com/Ikit24/pokedexcli/internal/pokeapi"
    "github.com/Ikit24/pokedexcli/internal/pokecache"
)

type Config struct {
    Next     string
    Previous string
    MyMap    map[string]CliCommand
    Cache    pokecache.Cache
    Caught   map[string]pokeapi.BattlePokemon
    Battle   map[string]pokeapi.BattlePokemon
}

type CliCommand struct {
    Name        string
    Description string
    Callback    func(*Config, []string) error
}
