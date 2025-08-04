package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "log"
    "io"
    "net/http"
)

func startRepl() {
    usr_input := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        usr_input.Scan()
        words := cleanInput(usr_input.Text())
        if len(words) == 0 {
            continue
        }
        commandName := words[0]

        cmd, ok := getCommands()[commandName]
        if ok {
            err := cmd.callback()
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
	callback    func() error
}

func getCommands() map[string]cliCommand {
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
    }
}
