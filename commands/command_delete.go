package commands

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"github.com/Ikit24/pokedexcli/internal/config"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)

func CommandDelete(cfg *config.Config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("delete command doesn't require any arguments.")
	}
	fmt.Println("WARNING! You are about to delete your current save. Are you sure? (y/n):")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("invalid command. Please type 'y' or 'n'.")
	}
	response := strings.TrimSpace(strings.ToLower(choice))
	if response == "n" {
		fmt.Println("Save deletion aborted.")
		return nil
	} else if response == "y" {
		fmt.Println("Deleting save file...")
		err = os.Remove("pokedex.json")
		if err != nil {
			return fmt.Errorf("save deletion was unsuccessfull. Please try again later.")
		}
		cfg.Caught = make(map[string]pokeapi.BattlePokemon)
		fmt.Println("Save deleted successfully.")
	}
	return nil
}
