package main

import (
    "fmt"
    "strings"
    "github.com/Ikit24/pokedexcli/internal/pokecache"
    "github.com/Ikit24/pokedexcli/internal/pokeapi"
    "time"
    "github.com/eiannone/keyboard"
)

func startRepl() {
    err := keyboard.Open()
    if err != nil {
        panic(err)
    }
    defer keyboard.Close()

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

    fmt.Print("Pokedex > ")
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
                    fmt.Println()
                    continue
                }
                fmt.Print("\r\033[K")

                commandName := words[0]
                cmd, ok := cfg.MyMap[commandName]
                if ok {
                    err := cmd.callback(&cfg, words[1:])
                    if err != nil {
                        fmt.Println(err)
                    }
                } else {
                    fmt.Println("Unknown command")
                }
                fmt.Println()

                commandHistory = append(commandHistory, commandText)
                // Clear input_buffer
                input_buffer = []rune{}
                // Reset to new input state
                historyIndex = len(commandHistory)
            } else if key == keyboard.KeyArrowUp {
                if historyIndex == 0 {
                    historyIndex = 0
                } else {
                    historyIndex -= 1
                    input_buffer = []rune(commandHistory[historyIndex])
                }
            } else if key == keyboard.KeyArrowDown {
                if historyIndex < len(commandHistory) {
                    historyIndex += 1
                }
                if historyIndex < len(commandHistory) {
                    input_buffer = []rune(commandHistory[historyIndex])
                } else {
                    input_buffer = []rune{}
                }
            } else if key == keyboard.KeySpace {
                input_buffer = append(input_buffer, ' ')
            } else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
                if len(input_buffer) > 0 {
                    input_buffer = input_buffer[:len(input_buffer)-1]
                }
            }
        } else {
            input_buffer = append(input_buffer, char)
        }
        fmt.Print("\r\033[K")
        fmt.Printf("Pokedex > %s", string(input_buffer))
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
            description: "Displays commamd list",
            callback:    commandHelp,
        },
        "exit": {
            name:        "exit",
            description: "Exits the Pokedex",
            callback:    commandExit,
        },
        "map": {
            name:        "map",
            description: "Displays the map to explore",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays the previous map list",
            callback:    commandMapb,
        },
        "explore": {
            name:        "explore",
            description: "Displays the pokemons in the current area",
            callback:    commandExplore,
        },
        "catch": {
            name:        "catch",
            description: "Catches Pokemons and adds them to the your Pokedex",
            callback:    commandCatch,
        },
        "inspect": {
            name:        "inspect",
            description: "Displays information about the already captured Pokemons",
            callback:    commandInspect,
        },
        "pokedex": {
            name:        "pokedex",
            description: "List of all the Pokemons that you captured",
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
