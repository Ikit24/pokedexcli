package main

import (
	"fmt"
)

func commandPokedex(cfg *config, _ []string) error {
	fmt.Println("Your Pokedex:")
	for name := range cfg.Caught {
		fmt.Println(" - " + name)
	}
	return nil
}
