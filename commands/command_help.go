package commands

import (
	"fmt"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandHelp(cfg *config.Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range cfg.MyMap {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
