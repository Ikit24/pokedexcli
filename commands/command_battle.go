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

func playerTurn(playerBattlePokemon, opponentBattlePokemon *pokeapi.BattlePokemon, playerMoves []Move) error {
	fmt.Println("Choose your move:")

	playerAttack, err := getStatValue(playerBattlePokemon.Stats, "attack")
	if err != nil {
		return fmt.Errorf("Error getting player attack")
	}
	opponentDefense, err := getStatValue(opponentBattlePokemon.Stats, "defense")
	if err != nil {
		return fmt.Errorf("Error getting opponent defense")
	}

	for i, move := range playerMoves {
		fmt.Printf("%d. %s\n", i+1, move.Name)
	}
	playerInput := bufio.NewReader(os.Stdin)
	choice, err := playerInput.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Invalid input")
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
	damage := (playerAttack * chosenMove.Power) / (opponentDefense * 2)
	finalDamage := int(damage)
	opponentBattlePokemon.CurrentHP -= finalDamage
	if opponentBattlePokemon.CurrentHP < 0 {
		opponentBattlePokemon.CurrentHP = 0
	}
	fmt.Printf("\n%s dealt %d to %s, remaining opponent HP: %d \n", playerBattlePokemon.Name,
		finalDamage,
		opponentBattlePokemon.Name,
		opponentBattlePokemon.CurrentHP)
	return nil
}

func opponentTurn(opponentBattlePokemon, playerBattlePokemon *pokeapi.BattlePokemon, opponentMoves []Move) error {
	opponentAttack, err := getStatValue(opponentBattlePokemon.Stats, "attack")
	if err != nil {
		return fmt.Errorf("Error getting opponent attack." )
	}
	playerDefense, err := getStatValue(playerBattlePokemon.Stats, "defense")
	if err != nil {
		return fmt.Errorf("Error getting player defense.")
	}
	
	strongestOpponentMove := opponentMoves[0]
	for _, move := range opponentMoves {
		if move.Power > strongestOpponentMove.Power {
			strongestOpponentMove = move
		}
	}
	opponentDamage := (opponentAttack * strongestOpponentMove.Power) / (playerDefense * 2)
	finalOpponentDamage := int(opponentDamage)
	playerBattlePokemon.CurrentHP -= finalOpponentDamage
	if playerBattlePokemon.CurrentHP < 0 {
		playerBattlePokemon.CurrentHP = 0
	}
	fmt.Printf("%s dealt %d to %s, remaining HP until %s faints: %d ", opponentBattlePokemon.Name,
	finalOpponentDamage,
	playerBattlePokemon.Name,
	playerBattlePokemon.Name,
	playerBattlePokemon.CurrentHP)
	fmt.Println()
	return nil
}


func checkVictory(cfg *config.Config, opponentBattlePokemon *pokeapi.BattlePokemon, playerBattlePokemon *pokeapi.BattlePokemon) error {
	fmt.Printf("You won! If you wish now is the time to capture %s without the chance of them escaping! ", opponentBattlePokemon.Name)
	playerBattlePokemon.CurrentXP += opponentBattlePokemon.BaseExperience

	fmt.Printf("%s gained %d XP!\n", playerBattlePokemon.Name, opponentBattlePokemon.BaseExperience)
	checkLevelUp(playerBattlePokemon)

	fmt.Printf("Would you like to capture %s? (y/n):\n", opponentBattlePokemon.Name)
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Invalid command. Please type 'y' or 'n'")
	}
	response := strings.TrimSpace(strings.ToLower(choice))
	if response == "n" {
		fmt.Printf("The defeated %s stays free. You won and walk on your path.\n", opponentBattlePokemon.Name)
	} else if response != "y" {
		return fmt.Errorf("Invalid response. Please enter y or n")
	} else {
		cfg.Caught[opponentBattlePokemon.Name] = *opponentBattlePokemon
		fmt.Printf("You caught %s!\n", opponentBattlePokemon.Name)
		return nil
	}
	return nil
}

func CommandBattle(cfg *config.Config, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: battle <your-pokemon> <target-pokemon>")
	}
	pokemonName := args[0]
	targetName := args[1]

	pokemon, ok := cfg.Caught[pokemonName]
	if !ok {
		return fmt.Errorf("Pokemon does not exist in your collection yet.")
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
	fmt.Println()
	fmt.Println("Proceed with battle? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("Invalid command")
	}
	response := strings.TrimSpace(strings.ToLower(line))
	if response == "n" || response == "no" {
		fmt.Println("Battle cancelled.")
		return nil
	}
	if response != "y" && response != "yes" {
		return fmt.Errorf("Invalid response. Please enter y or n")
	}
	fmt.Println("Battle begins!")
	fmt.Println()
	
	// Create battle copies
	playerBattlePokemon := cfg.Caught[pokemonName]
	opponentBattlePokemon := cfg.Battle[targetName]

	playerMoves, err := generateBasicMoves(playerBattlePokemon)
	if err != nil {
		return fmt.Errorf("Error generating player moves: %w", err)
	}
	opponentMoves, err := generateBasicMoves(opponentBattlePokemon)
	if err != nil {
		return fmt.Errorf("Error generating opponent moves: %w", err)
	}

	playerMaxHP, err := getStatValue(playerBattlePokemon.Stats, "hp")
	if err != nil {
		return fmt.Errorf("Error, couldn't initialize player hp. Please try again.")
	}
	playerBattlePokemon.CurrentHP = playerMaxHP
	playerBattlePokemon.StatusEffects = []string{}
	playerBattlePokemon.CurrentXP = 0
	playerBattlePokemon.Level = 1

	opponentMaxHP, err := getStatValue(opponentBattlePokemon.Stats, "hp")
	if err != nil {
		return fmt.Errorf("Error, couldn't initialize opponent hp. Please try again.")
	}
	opponentBattlePokemon.CurrentHP = opponentMaxHP
	opponentBattlePokemon.StatusEffects = []string{}
	opponentBattlePokemon.CurrentXP = 0
	opponentBattlePokemon.Level = 1

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
	} else {
		fmt.Println("Opponent is faster. You go second!")
	}

	for playerBattlePokemon.CurrentHP > 0 && opponentBattlePokemon.CurrentHP > 0 {
		if playerSpeed >= opponentSpeed {
			err := playerTurn(&playerBattlePokemon, &opponentBattlePokemon, playerMoves)
			if err != nil {
				return err
			}
			if opponentBattlePokemon.CurrentHP <= 0 {
				err := checkVictory(cfg, &opponentBattlePokemon, &playerBattlePokemon)
				if err != nil {
					return err
				}
				return nil
				err = opponentTurn(&opponentBattlePokemon, &playerBattlePokemon, opponentMoves)
				if err != nil {
						return err
				}
				if playerBattlePokemon.CurrentHP <= 0 {
					fmt.Println("You lost... You walk away defeated")
					return nil
				}
			}
		} else {
			err := opponentTurn(&opponentBattlePokemon, &playerBattlePokemon, opponentMoves)
			if err != nil {
				return err
			}
			if playerBattlePokemon.CurrentHP <= 0 {
				fmt.Println("You lost... You walk away defeated")
				return nil
			}
			err = playerTurn(&playerBattlePokemon, &opponentBattlePokemon, playerMoves)
			if err != nil {
				return err
			}
			if playerBattlePokemon.CurrentHP <= 0 {
			fmt.Println("You lost... You walk away defeated")
			return nil
			}
		
			if opponentBattlePokemon.CurrentHP <= 0 {
				err := checkVictory(cfg, &opponentBattlePokemon, &playerBattlePokemon)
				if err != nil {
					return err
				}			
				return nil
			}
		}
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
		return nil, fmt.Errorf("Cannot fetch stats")
	}
    moves = append(moves, Move{
        Name:     "Tackle",
        Power:    30 + (attackStat / 10),
        Accuracy: 95,
    })
	if len(pokemon.Types) > 0 {
		moves = append(moves, getTypeMove(pokemon.Types[0].Type.Name))
	} else {
		return nil, fmt.Errorf("Pokemon %s has no type data - possible corruption. Please restart.", pokemon.Name)
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
    return 0, fmt.Errorf("Pokemon %s stat not found, data corruption or server error possible. Please restart.", statName)
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

func getXPForLevel(level int) int {
	return level * level * level
}

func getLevelFromXP(currentXP int) int {
	level := 1
	for getXPForLevel(level+1) <= currentXP {
		level++
	}
	return level
}

func checkLevelUp(pokemon *pokeapi.BattlePokemon) bool {
	newLevel := getLevelFromXP(pokemon.CurrentXP)
	if newLevel > pokemon.Level {
		fmt.Printf("%s leveled up! Now level %d\n", pokemon.Name, newLevel)
		pokemon.Level = newLevel
		return true
	}
	return false
}
