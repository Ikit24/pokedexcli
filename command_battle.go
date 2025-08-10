package main

import (
	"fmt"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)

func commandBattle(cfg *config, pokemon_name, target_name string) error {
	if len(pokemon_name) == 0 {
		return fmt.Errorf("Please provide pokemon name in order to battle")
	}
	if len(target_name) == 0 {
		return fmt.Errorf("Please provide a target pokemon name in your area in order to battle")
	}
}

type Move struct {
    Name     string
    Power    int
    Accuracy int
}

func generateBasicMoves(pokemon pokeapi.BattlePokemon) ([]Move, error) {
    moves := []Move{}
 
    attackStat, err := getStatValue(pokemon.Stats, "attack")
	if err != nil {
		return nil, fmt.Errorf("cannot fetch stats")
	}
    moves = append(moves, Move{
        Name:     "Tackle",
        Power:    30 + (attackStat / 10),
        Accuracy: 95,
    })
	if len(pokemon.Types) > 0 {
		moves = append(moves, getTypeMove(pokemon.Types[0].Type.Name))
	} else {
		return nil, fmt.Errorf("pokemon %s has no type data - possible corruption. Please restart.", pokemon.Name)
	}

    return moves, nil
}

func getStatValue(stats []struct {
	BaseStat int `json:"base_stat"`
	Stat	 struct {
		Name string `json:"name"`
    } `json:"stat"`
}, statName string) (int, error) {
    for _, stat := range stats {
        if stat.Stat.Name == statName {
            return stat.BaseStat
        }
    }
    return 0, fmt.Errorf("pokemon %s stat not found, data corruption or server error possible. Please restart.", statName)
}

func getTypeMove(pokemonType string) Move {
	switch pokemonType {
	case "electric":
		return Move{Name: "Thunder Shock", Power: 40, Accuracy: 100}
	case "fire":
		return Move{Name: "Ember", Power: 40, Accuracy: 100}
	case "water":
		return Move{Name: "Water Gun", Power: 40, Accuracy: 100}
	default:
		return Move{Name: "Normal Attack", Power: 35, Accuracy: 100}
	}
}
