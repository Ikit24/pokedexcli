package main

import (
    "fmt"
    "strings"
    "github.com/Ikit24/pokedexcli/internal/config"
    "github.com/Ikit24/pokedexcli/commands"
    "github.com/eiannone/keyboard"
)

func startRepl(cfg *config.Config) {
    err := keyboard.Open()
    if err != nil {
        panic(err)
    }
    defer keyboard.Close()

    var input_buffer []rune
    var commandHistory []string
    var historyIndex int

    historyIndex = len(commandHistory)

    fmt.Print("Welcome to Pokdexcli! If you don't know how to start type 'help'\n")
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
                    err := cmd.Callback(cfg, words[1:])
                    if err != nil {
                        fmt.Println(err)
                    }
                } else {
                    fmt.Println("unknown command")
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

func getCommands(cfg *config.Config) map[string]config.CliCommand {
    return map[string]config.CliCommand{
        "help": {
            Name:        "- help",
            Description: "Displays the command list.",
            Callback:    commands.CommandHelp,
        },
        "exit": {
            Name:        "- exit",
            Description: "Exits the Pokedex.",
            Callback:    commands.CommandExit,
        },
        "map": {
            Name:        "- map",
            Description: "Displays the map of the current area.",
            Callback:    commands.CommandMap,
        },
        "mapb": {
            Name:        "- mapb",
            Description: "Displays the map of the previous area.",
            Callback:    commands.CommandMapb,
        },
        "explore": {
            Name:        "- explore",
            Description: "Displays the pokemons in the currently explored area.",
            Callback:    commands.CommandExplore,
        },
        "catch": {
            Name:        "- catch",
            Description: "Catches Pokemons and adds them to your Pokedex.",
            Callback:    commands.CommandCatch,
        },
        "inspect": {
            Name:        "- inspect",
            Description: "Displays information about the Pokemon in your Pokedex.",
            Callback:    commands.CommandInspect,
        },
        "pokedex": {
            Name:        "- pokedex",
            Description: "List of all the Pokemons that you captured.",
            Callback:    commands.CommandPokedex,
        },
        "battle": {
            Name:        "- battle",
            Description: "Initiates a battle by using one of your captured Pokemon and a selected pokemon in the current area.",
            Callback:    commands.CommandBattle,
        },
        "save": {
            Name:        "- save",
            Description: "Saves current progress.",
            Callback:    commands.CommandSave,
        },
        "delete": {
            Name:        "- delete",
            Description: "Deletes current save file. WARNING! This cannot be undone after executed.",
            Callback:    commands.CommandDelete,
        },
        "evolve": {
            Name:        "- evolve",
            Description: "Optional evolve check to see if any of your Pokemons evolved outside of battles.",
            Callback:    commands.CommandEvolve,
        },
    }
}
