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

func CommandSave(cfg *config.Config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("Usage: save (without arguments)")
	}
	save, err := json.Marshal(cfg.Caught)
	if err != nil {
		return fmt.Errorf("Error, marshal to JSON failed!")
	}
	err = os.WriteFile("pokedex.json", save, 0644)
	if err != nil {
		fmt.Println("Save failed! Retry? (y/n):")
		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Invalid command. Please type 'y' or 'n'")
		}
		response := strings.TrimSpace(strings.ToLower(choice))
		if response == "n" {
			fmt.Println("Save unsuccessfull... Continuing...")
			return nil
		} else if response == "y" {
			fmt.Println("Retrying...")
			err = os.WriteFile("pokedex.json", save, 0644)
			if err != nil {
				fmt.Println("Save failed again. Continuing without saving.")
				fmt.Println("Please try saving later.")
			} else {
				fmt.Println("Save successfull!")
			}
		} else {
			fmt.Println("Invlaid response. Save cancelled, continuing...")
		}
	} else {
		fmt.Println("Save successfull!")
		return nil
	}
	return nil
}
