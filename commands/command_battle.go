package commands

import (
	"fmt"
	"bufio"
	"encoding/json"
	"strings"
	"os"
	"strconv"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/config"
	"math/rand"
	"context"
	"net/http"
	"time"
)

func playerTurn(playerBattlePokemon, opponentBattlePokemon *pokeapi.BattlePokemon, playerMoves []Move) error {
	playerAttack, err := getStatValue(playerBattlePokemon.Stats, "attack")
	if err != nil {
		return fmt.Errorf("Error getting player attack")
	}
	opponentDefense, err := getStatValue(opponentBattlePokemon.Stats, "defense")
	if err != nil {
		return fmt.Errorf("Error getting opponent defense")
	}

	var chosenMove Move

	for {
		fmt.Println("\nChoose your move:")
		for i, move := range playerMoves {
			fmt.Printf("%d. %s\n", i+1, move.Name)
		}

		playerInput := bufio.NewReader(os.Stdin)
		choice, err := playerInput.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		response := strings.TrimSpace(strings.ToLower(choice))
		if len(response) == 0 {
			fmt.Println("Please enter a number.")
			continue
		}

		choiceNum, err := strconv.Atoi(response)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}

		if choiceNum < 1 || choiceNum > len(playerMoves) {
			fmt.Printf("Please enter the number between 1 and  %d.\n", len(playerMoves))
			continue
		}
		chosenMove = playerMoves[choiceNum - 1]
		break
	}

	damage := (playerAttack * chosenMove.Power) / (opponentDefense * 2)
	finalDamage := int(damage)
	opponentBattlePokemon.CurrentHP -= finalDamage
	if opponentBattlePokemon.CurrentHP < 0 {
		opponentBattlePokemon.CurrentHP = 0
	}
	fmt.Printf("\n%s dealt %s to %s, remaining opponent HP: %s \n", 
		colorize("\033[32m", playerBattlePokemon.Name),
		colorize("\033[31m", strconv.Itoa(finalDamage)),
		colorize("\033[35m", opponentBattlePokemon.Name),
		colorize("\033[34m", strconv.Itoa(opponentBattlePokemon.CurrentHP)))
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
	fmt.Printf("\n%s dealt %s to %s, remaining HP until %s faints: %s ",
		colorize("\033[35m", opponentBattlePokemon.Name),
		colorize("\033[31m", strconv.Itoa(finalOpponentDamage)),
		colorize("\033[32m", playerBattlePokemon.Name),
		colorize("\033[32m", playerBattlePokemon.Name),
		colorize("\033[33m", strconv.Itoa(playerBattlePokemon.CurrentHP)))
	fmt.Println()
	return nil
}

func checkVictory(cfg *config.Config, opponentBattlePokemon *pokeapi.BattlePokemon, playerBattlePokemon *pokeapi.BattlePokemon) error {
    fmt.Printf("You won! If you wish now is the time to capture %s without the chance of them escaping! ", colorize("\033[35m", opponentBattlePokemon.Name))
    playerBattlePokemon.CurrentXP += opponentBattlePokemon.BaseExperience

    fmt.Printf("%s gained %s XP!\n",
		colorize("\033[32m", playerBattlePokemon.Name),
		colorize("\033[36m", strconv.Itoa(opponentBattlePokemon.BaseExperience)))
    checkLevelUp(playerBattlePokemon)
	givePartyXP(cfg, opponentBattlePokemon.BaseExperience)

	cfg.Caught[playerBattlePokemon.Name] = *playerBattlePokemon

	fmt.Printf("Would you like to capture %s? (y/n):\n", colorize("\033[35m", opponentBattlePokemon.Name))
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Invalid command. Please type 'y' or 'n'")
	}
	response := strings.TrimSpace(strings.ToLower(choice))
	if response == "n" {
		fmt.Printf("The defeated %s stays free.\n", colorize("\033[35m", opponentBattlePokemon.Name))
		AutoSave(cfg)
	} else if response != "y" {
		return fmt.Errorf("Invalid response. Please enter y or n")
	} else {
		caught := *opponentBattlePokemon
		caught.CaughtAt = time.Now().UTC()
		caught.EvolutionDelaySecs = 3600
		caught.HasEvolved = false

		next, min, err := DetermineNextEvolution(cfg, caught.Name)
		if err != nil {
			fmt.Printf("Warning: could not determine evolution for %s: %v\n", caught.Name, err)
		}
		caught.EvolvesTo = next
		caught.MinLevelForEvolution = min

		cfg.Caught[caught.Name] = caught
		fmt.Printf("You caught %s!\n", colorize("\033[35m", opponentBattlePokemon.Name))

		evolved, err := RunEvolutionPass(context.Background(), http.DefaultClient, cfg.Caught)
		if err != nil {
			fmt.Printf("Evolution check failed: %v\n", err)
		}
		for _, msg := range evolved {
			fmt.Println(msg)
		}
		AutoSave(cfg)
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
	for {
		fmt.Println("Proceed with battle? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		response := strings.TrimSpace(strings.ToLower(line))

		if response == "n" || response == "no" {
			fmt.Println("Battle cancelled.")
			return nil
		}
		if response == "y" || response == "yes" {
			break
		}
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
	if playerBattlePokemon.CurrentXP == 0 && playerBattlePokemon.Level == 0 {
		playerBattlePokemon.Level = 1
	}

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
			}
			if opponentBattlePokemon.CurrentHP > 0 {
				err = opponentTurn(&opponentBattlePokemon, &playerBattlePokemon, opponentMoves)
				if err != nil {
					return err
				}
				if playerBattlePokemon.CurrentHP <= 0 {
					fmt.Println("You lost... You walk away defeated")
					AutoSave(cfg)
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
				AutoSave(cfg)
				return nil
			}
			err = playerTurn(&playerBattlePokemon, &opponentBattlePokemon, playerMoves)
			if err != nil {
				return err
			}
			if playerBattlePokemon.CurrentHP <= 0 {
				fmt.Println("You lost... You walk away defeated")
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

func GetXPForLevel(level int) int {
	return level * level * level
}

func GetLevelFromXP(currentXP int) int {
	level := 1
	for GetXPForLevel(level+1) <= currentXP && level < 100 {
		level++
	}
	return level
}

func checkLevelUp(pokemon *pokeapi.BattlePokemon) bool {
	if pokemon.Level >= 100 {
		fmt.Printf("%s is already at max level (100)!\n", colorize("\033[32m", pokemon.Name))
		return false
	}

	newLevel := GetLevelFromXP(pokemon.CurrentXP)
	if newLevel > pokemon.Level {
		if newLevel > 100 {
			newLevel = 100
		}

		fmt.Printf("%s leveled up! Now level %s\n",
			colorize("\033[32m", pokemon.Name),
			colorize("\033[36m", strconv.Itoa(newLevel)))
		
		levelUpStats(pokemon)

		if newLevel == 100 {
			fmt.Printf("%s has reached the maximum level!\n", colorize("\033[32m", pokemon.Name))
		}
		pokemon.Level = newLevel
		return true
	}
	return false
}

func AutoSave(cfg *config.Config) error {
	save, err := json.Marshal(cfg.Caught)
	if err != nil {
		return fmt.Errorf("Error, marshal to JSON failed!")
	}
	err = os.WriteFile("pokedex.json", save, 0644)
	if err != nil {
		fmt.Println("Autosave failed. Retrying...")
		err = os.WriteFile("pokedex.json", save, 0644)
		if err != nil {
			fmt.Printf("Warning: Could not save progress: %v\n", err)
		}
	}
	return nil
}

func colorize(color, text string) string {
	return color + text + "\033[0m"
}

func levelUpStats(pokemon *pokeapi.BattlePokemon) {
	for i := range pokemon.Stats{
		increase := pokemon.Stats[i].BaseStat * (rand.Intn(2) + 2) / 100
		pokemon.Stats[i].BaseStat += increase
	}
}

func givePartyXP(cfg *config.Config, baseXP int) {
	for name, pokemon := range cfg.Caught {
		pokemon.CurrentXP += baseXP / 5 
		checkLevelUp(&pokemon)
		cfg.Caught[name] = pokemon
	}
}
