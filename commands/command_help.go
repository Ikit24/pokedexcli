package commands

import (
	"fmt"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandHelp(cfg *config.Config, args []string) error {
    if len(args) > 0 {
        return fmt.Errorf("help command doesn't require any arguments.")
    }
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range cfg.MyMap {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
