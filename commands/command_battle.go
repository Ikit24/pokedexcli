package commands

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"strconv"
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
	
	// Create battle copies
	playerBattlePokemon := cfg.Caught[pokemonName]
	opponentBattlePokemon := cfg.Battle[targetName]

	playerMoves, err := generateBasicMoves(playerBattlePokemon)
	if err != nil {
		return fmt.Errorf("Error generating player moves: %w", err)
	}

	playerMaxHP, err := getStatValue(playerBattlePokemon.Stats, "hp")
	if err != nil {
		return fmt.Errorf("Error, couldn't initialize player hp. Please try again.")
	}
	playerBattlePokemon.CurrentHP = playerMaxHP
	playerBattlePokemon.StatusEffects = []string{}

	opponentMaxHP, err := getStatValue(opponentBattlePokemon.Stats, "hp")
	if err != nil {
		return fmt.Errorf("Error, couldn't initialize opponent hp. Please try again.")
	}
	opponentBattlePokemon.CurrentHP = opponentMaxHP
	opponentBattlePokemon.StatusEffects = []string{}

	playerSpeed, err := getStatValue(playerBattlePokemon.Stats, "speed")
	if err != nil {
		return fmt.Errorf("Error getting player speed")
	}

	opponentSpeed, err := getStatValue(opponentBattlePokemon.Stats, "speed")
	if err != nil {
		return fmt.Errorf("Error getting opponent speed")
	}

	if playerSpeed >= opponentSpeed {
		fmt.Println("Your Pokemon is faster. You go first!")
		fmt.Println("Choose your move:")
		for i, move := range playerMoves {
			fmt.Printf("%d. %s\n", i+1, move.Name)
		}
		playerInput := bufio.NewReader(os.Stdin)
		choice, err := playerInput.ReadString('\n')
		if err != nil {
			return fmt.Errorf("invalid input")
		}
		response := strings.TrimSpace(strings.ToLower(choice))

		if len(response) == 0 {
			return fmt.Errorf("Please enter a number")
		}
		choiceNum, err := strconv.Atoi(response)
		if err != nil {
			return fmt.Errorf("Please enter a number")
		} else if choiceNum < 1 || choiceNum > len(playerMoves) {
			return fmt.Errorf("Please enter the number between 1 and  %d", len(playerMoves))
		}
		
		chosenMove := playerMoves[choiceNum - 1]
		playerAttack, err := getStatValue(playerBattlePokemon.Stats, "attack")
		if err != nil {
			return fmt.Errorf("Error getting player attack.")
		}
		playerDefense, err := getStatValue(playerBattlePokemon.Stats, "defense")
		if err != nil {
			return fmt.Errorf("Error getting player defense.")
		}

		opponentAttack, err := getStatValue(opponentBattlePokemon.Stats, "attack")
		if err != nil {
			return fmt.Errorf("Error getting opponent attack.")
		}
		opponentDefense, err := getStatValue(opponentBattlePokemon.Stats, "defense")
		if err != nil {
			return fmt.Errorf("Error getting opponent defense.")
		}
		damage := (playerAttack * chosenMove.Power) / (opponentDefense * 2)
		finalDamage := int(damage)
		opponentBattlePokemon.CurrentHP -= finalDamage
		fmt.Printf("%s dealt %d to %s, remaining opponent HP: %d", playerBattlePokemon.Name,
			finalDamage,
			opponentBattlePokemon.Name,
			opponentBattlePokemon.CurrentHP)

		if opponentBattlePokemon.CurrentHP <= 0 {
			fmt.Println("You won! If you wish you can catpure the opponent with 100% success.")
			reader := bufio.NewReader(os.Stdin)
			choice, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("Invalid command. Please type 'y' or 'n'")
			}
			response := strings.TrimSpace(strings.ToLower(choice))
			if response == "n" {
				fmt.Println("The defeated pokemon stays free. You won and walk on your path.")
			} else if response != "y" {
				return fmt.Errorf("invalid response. Please enter y or n")
			} else {
				cfg.Caught[opponentBattlePokemon.Name] = opponentBattlePokemon
				fmt.Printf("You caught %s!\n", opponentBattlePokemon.Name)
			}
		}
	} else {
		fmt.Println("Opponent is faster. You go second!")
	}

	opponentMoves, err := generateBasicMoves(opponentBattlePokemon)
	if err != nil {
		return fmt.Errorf("Error generating opponent moves: %w", err)
	}
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
