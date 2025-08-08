package main

import (
	"fmt"
)

func commandInspect(cfg *config, pokemon_name []string) error {
	if len(pokemon_name) == 0 {
		return fmt.Errorf("Must provide pokemon name in order to display it's details")
	}
	p, ok := cfg.Caught[pokemon_name[0]]
	if ok == false {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t:= range p.Types {
		fmt.Println(" - " + t.Type.Name)
	}
	return nil
}
