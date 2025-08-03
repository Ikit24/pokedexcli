package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}


func startRepl() {
    usr_input := bufio.NewScanner(os.Stdin)

    commands := map[string]cliCommand{}

    commandHelp := func() error {
        fmt.Println("Welcome to the Pokedex!\nUsage:")
        for name, cmd := range commands {
            fmt.Printf("%s: %s\n", name, cmd.description)
        }
        return nil
    }

    commands["help"] = cliCommand{
            name:           "help",
            description:    "Displays a help message",
            callback:       commandHelp,
    }
    commands["exit"] = cliCommand{
            name:           "exit",
            description:    "Exit the Pokedex",
            callback:       commandExit,
    }


    for {
        fmt.Print("Pokedex > ")
        usr_input.Scan()
        words := cleanInput(usr_input.Text())
        if len(words) == 0 {
            continue
        }
        commandName := words[0]

        cmd, ok := commands[commandName]
        if ok {
            err := cmd.callback()
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}


func cleanInput(text string) []string {
    output := strings.ToLower(text)
    words := strings.Fields(output)
    return words
}


func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}
