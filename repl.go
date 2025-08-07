package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "github.com/Ikit24/pokedexcli/internal/pokecache"
    "time"
)

func startRepl() {
    usr_input := bufio.NewScanner(os.Stdin)

    var cfg config

    cfg.Next = "https://pokeapi.co/api/v2/location-area/"
    cfg.Previous = ""

    cfg.Cache = pokecache.NewCache(1 * time.Minute)

    cfg.MyMap = getCommands(&cfg)

    for {
        fmt.Println("Pokedex > ")
        usr_input.Scan()
        words := cleanInput(usr_input.Text())
        if len(words) == 0 {
            continue
        }
        commandName := words[0]

        cmd, ok := cfg.MyMap[commandName]
        if ok {
            err := cmd.callback(&cfg, words[1:])
            if err != nil {
                fmt.Println(err)
            }
            continue
        } else {
            fmt.Println("Unknown command")
            continue
        }
    }
}

func cleanInput(text string) []string {
    output := strings.ToLower(text)
    words := strings.Fields(output)
    return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands(cfg *config) map[string]cliCommand {
    return map[string]cliCommand{
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "map": {
            name:        "map",
            description: "Displays the current 20 entries",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays the previous 20 entries",
            callback:    commandMapb,
        },
        "explore": {
            name:        "explore",
            description: "Displays pokemons in the current area",
            callback:    commandExplore,
        },
        "catch": {
            name:        "catch",
            description: "Catches Pokemon adds them to the user's Pokedex",
            callback:    commandCatch,
        },
    }
}

type config struct {
    Next        string
    Previous    string
    MyMap       map[string]cliCommand
    Cache       pokecache.Cache
}
