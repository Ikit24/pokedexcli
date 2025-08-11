package commands

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandBattle(cfg *config.Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: battle <your-pokemon> <target-pokemon>")
	}
	pokemonName := args[0]
	targetName := args[1]

	pokemon, ok := cfg.Caught[pokemonName]
	if !ok {
		return fmt.Errorf("Pokemon does not exist in your collection")
	}
	targetPokemon, ok := cfg.Battle[targetName]
	if !ok {
		return fmt.Errorf("Invalid target")
	}

	playerPokemonStats, err := getAllBattleStats(pokemon)
	if err != nil {
		return err
	}

	opponentPokemonStats, err := getAllBattleStats(targetPokemon)
	if err != nil {
		return err
	}

	displayBattleComparison(playerPokemonStats, opponentPokemonStats, pokemon.Name, targetPokemon.Name)
	fmt.Println("Proceed with battle? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("invalid command")
	}
	response := strings.TrimSpace(strings.ToLower(line))
	if response == "n" || response == "no" {
		fmt.Println("Battle cancelled.")
		return nil
	}
	if response != "y" && response != "yes" {
		return fmt.Errorf("invalid response. Please enter y or n")
	}
	fmt.Println("Battle begins!")
	return nil
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
            return stat.BaseStat, nil
        }
    }
    return 0, fmt.Errorf("pokemon %s stat not found, data corruption or server error possible. Please restart.", statName)
}

func getAllBattleStats(pokemon pokeapi.BattlePokemon) (map[string]int, error) {
	stats := make(map[string]int)
	for _, stat := range pokemon.Stats {
		stats[stat.Stat.Name] = stat.BaseStat
	}
	return stats, nil
}

func displayBattleComparison(playerStats, opponentStats map[string]int, playerName, opponentName string) {
	fmt.Printf("%-20s vs %-20s\n", playerName, opponentName)
	fmt.Println(strings.Repeat("-", 45))

	coreStats := []string{"hp", "attack", "defense", "speed"}

	for _, statName := range coreStats {
		playerValue := playerStats[statName]
		opponentValue := opponentStats[statName]

		fmt.Printf("%-20s %-20s\n", 
            fmt.Sprintf("%s: %d", strings.Title(statName), playerValue),
            fmt.Sprintf("%s: %d", strings.Title(statName), opponentValue))
	}
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
