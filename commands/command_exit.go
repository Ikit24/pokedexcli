package commands

import (
	"fmt"
	"os"
    "github.com/Ikit24/pokedexcli/internal/config"
)

func CommandExit(cfg *config.Config, _ []string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}
