package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
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
        fmt.Printf("Your command was: %s\n", commandName)
    }
}

func cleanInput(text string) []string {
    output := strings.ToLower(text)
    words := strings.Fields(output)
    return words
}
