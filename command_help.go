package main

import (
	"fmt"
)

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for name, cmd := range cfg.MyMap {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}
