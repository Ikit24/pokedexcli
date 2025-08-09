package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "github.com/Ikit24/pokedexcli/internal/pokecache"
    "github.com/Ikit24/pokedexcli/internal/pokeapi"
    "time"
    "github.com/eiannone/keyboard"
)

func startRepl() {
    kb, err := keyboard.Open()
    if err != nil {
        panic(err)
    }
    defer kb.Close()

    var input_buffer []rune
    var commandHistory []string
    var historyIndex int
    var cfg config

    historyIndex = len(commandHistory)

    cfg.Next = "https://pokeapi.co/api/v2/location-area/"
    cfg.Previous = ""
    cfg.Cache = pokecache.NewCache(1 * time.Minute)
    cfg.MyMap = getCommands(&cfg)
    cfg.Caught = make(map[string]pokeapi.Pokemon)

    for {
        char, key, err := keyboard.GetSingleKey()
        if err != nil {
            panic(err)
        }
        if key != 0 {
            if key == keyboard.KeyEnter {
                commandText := string(input_buffer)
                words := cleanInput(commandText)

                if len(words) == 0 {
                    input_buffer = []rune{}
                    historyIndex = len(commandHistory)
                    continue
                }
            commandName := words[0]
            cmd, ok := cfg.MyMap[commandName]
            if ok {
                err := cmd.callback(&cfg, words[1:])
                if err != nil {
                    fmt.Print("\r\033[K")
                    fmt.Println(err)
                }
            } else {
                fmt.Print("\r\033[K")
                fmt.Println("Unknown command")
            }
            commandHistory = append(commandHistory, commandText)
            // Clear input_buffer
            input_buffer = []rune{}
            // Reset to new input state
            historyIndex = len(commandHistory)

            } else if key == keyboard.KeyArrowDown {
                
            } else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
                if len(input_buffer) > 0 {
                    input_buffer = input_buffer[:len(input_buffer)-1]
                }
        } else {
            input_buffer = append(input_buffer, char)
        }
        fmt.Print("\r\033[K")
        fmt.Printf("Pokedex > %s", string(input_buffer))

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
        "inspect": {
            name:        "inspect",
            description: "Displays information about the already caputred Pokemon",
            callback:    commandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "list of all the names of the Pokemon the user has caught",
            callback:    commandPokedex,
        },
    }
}

type config struct {
    Next        string
    Previous    string
    MyMap       map[string]cliCommand
    Cache       pokecache.Cache
    Caught      map[string]pokeapi.Pokemon
}
