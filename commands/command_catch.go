package commands

import (
	"fmt"
	"bufio"
	"io"
	"os"
	"net/http"
	"encoding/json"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
	"github.com/Ikit24/pokedexcli/internal/config"
	"math/rand"
	"strconv"
	"strings"
)

const (
	baseXP_low		= 100
	baseXP_medium	= 150
	baseXP_high		= 300
)

func CommandCatch(cfg *config.Config, pokemonName []string) error {
	if len(pokemonName) == 0 {
		return fmt.Errorf("Must provide pokemon name in order to catch.")
	}
	
	var catch_URL = "https://pokeapi.co/api/v2/pokemon/" + pokemonName[0] + "/"

	var body []byte
	var err error

	cachedData, ok := cfg.Cache.Get(catch_URL)
	if !ok {
		res, httpErr := http.Get(catch_URL)
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(catch_URL, body)
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	} else {
		body = cachedData
		err = nil
	}
	var apiResponse pokeapi.BattlePokemon

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	validBalls := getValidPokeballs(apiResponse.BaseExperience)
	chosenBall := selectPokeball(validBalls)

	fmt.Printf("Throwing a %s at %s...\n", chosenBall.Name, pokemonName[0])
	fmt.Println()

	var catch_Chance int
	var baseCatchRate int

	if apiResponse.BaseExperience <= baseXP_low {
		baseCatchRate = 80
	} else if apiResponse.BaseExperience <= baseXP_medium {
		baseCatchRate = 50
	} else if apiResponse.BaseExperience <= baseXP_high {
		baseCatchRate = 20
	} else {
		baseCatchRate = 5
	}

	catch_Chance = int(float64(baseCatchRate) * (float64(chosenBall.CatchRate) / 100.0))
	if catch_Chance > 100 {
		catch_Chance = 100
	}

	randomRoll := rand.Intn(100) + 1
	if randomRoll <= catch_Chance {
		fmt.Printf("%s was caught!\n", pokemonName[0])
		if cfg.Caught == nil {
			cfg.Caught = make(map[string]pokeapi.BattlePokemon)
		}
		cfg.Caught[apiResponse.Name] = apiResponse
		AutoSave(cfg)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName[0])
	}
	return nil
}

type PokeBall struct {
	Name		string
	CatchRate	int
	MinimForAll bool
}

func getPokeBalls() []PokeBall {
	return []PokeBall{
		{Name: "Pokeball", CatchRate: 50, MinimForAll: false},
		{Name: "Great Ball", CatchRate: 75, MinimForAll: false},
		{Name: "Ultra Ball", CatchRate: 90, MinimForAll: false},
		{Name: "Master Ball", CatchRate: 100, MinimForAll: true},
	}
}

func getValidPokeballs(pokemonBaseExp int) []PokeBall {
	allBalls := getPokeBalls()
	validBalls := []PokeBall{}

	if pokemonBaseExp <= baseXP_medium {
		validBalls = allBalls
	} else if pokemonBaseExp <= baseXP_high {
		for _, ball := range allBalls {
			if ball.Name == "Great Ball" || ball.Name == "Ultra Ball" || ball.Name == "Master Ball" {
				validBalls = append(validBalls, ball)
			}
		}
	} else {
		for _, ball := range allBalls {
			if ball.Name == "Master Ball" {
				validBalls = append(validBalls, ball)
			}
		}
	}
	return validBalls
}

func selectPokeball(validBalls []PokeBall) PokeBall {
	var chosenBall PokeBall

	for {
		fmt.Println("\nChoose a Pokeball:")
		for i, ball := range validBalls {
			fmt.Printf("%d. %s\n", i+1, ball.Name)
		}

		playerInput :=bufio.NewReader(os.Stdin)
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

		if choiceNum < 1 || choiceNum > len(validBalls) {
			fmt.Printf("Please enter a number between 1 and %d.\n", len(validBalls))
			continue
		}
		chosenBall = validBalls[choiceNum - 1]
		break
	}
	return chosenBall
}
