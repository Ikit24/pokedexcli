package config

import (
    "github.com/Ikit24/pokedexcli/internal/pokeapi"
    "github.com/Ikit24/pokedexcli/internal/pokecache"
)

type Config struct {
    ExploredAreas []string
    MyMap         map[string]CliCommand
    Caught        map[string]pokeapi.BattlePokemon
    Battle        map[string]pokeapi.BattlePokemon
    Cache         pokecache.Cache
    Next          string
    Previous      string
}

type CliCommand struct {
    Name        string
    Description string
    Callback    func(*Config, []string) error
}
