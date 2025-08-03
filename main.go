package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    usr_input := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        usr_input.Scan()
        text := usr_input.Text()
        cleaned := strings.TrimSpace(text)
        split_input := strings.Split(cleaned, " ")
        if split_input[0] == "" {
            fmt.Print("Error empty string, please enter a string.")
            continue
        }
        fmt.Printf("Your command was: %s\n", split_input[0])
    }
}
