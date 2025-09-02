package commands

import (
	"fmt"
	"os"
    "strings"
    "bufio"
    "github.com/Ikit24/pokedexcli/internal/config"
)

func CommandExit(cfg *config.Config, args []string) error {
    if len(args) > 0 {
        return fmt.Errorf("exit command doesn't require any arguments.")
    }
    fmt.Println("You are about to exit. Would you like to save your progress first? (y/n) ")
    for {
        reader := bufio.NewReader(os.Stdin)
        choice, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Invalid command. Please type 'y' or 'n'.")
            continue
        }

        response := strings.TrimSpace(strings.ToLower(choice))
        if response == "n" {
            fmt.Println("Closing Pokedex without saving... Goodbye!")
            os.Exit(0)
        } else if response == "y" {
            fmt.Println("Saving data..")
            err = AutoSave(cfg)
            if err != nil {
                fmt.Println("Save failed %v. Please try again later.\n", err)
                return nil
            }
            fmt.Println("Save successful!")
            fmt.Println("Closing Pokedex... Goodbye!")
            os.Exit(0)
        } else {
            fmt.Println("Invalid command. Please type 'y' or 'n'.")
            continue
        }
    }
}
