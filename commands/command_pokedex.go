package commands

import (
	"fmt"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandPokedex(cfg *config.Config, _ []string) error {
	if len(cfg.Caught) == 0 {
		fmt.Println("Pokedex is empty.")
	} else {
		fmt.Println("Your Pokedex:")
		for name := range cfg.Caught {
			fmt.Println(" - " + name)
		}
	}
	return nil
}
