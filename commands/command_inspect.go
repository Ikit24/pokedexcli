package commands

import (
	"fmt"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func CommandInspect(cfg *config.Config, pokemon_name []string) error {
	if len(pokemon_name) == 0 {
		return fmt.Errorf("Must provide pokemon name in order to display it's details")
	}
	p, ok := cfg.Caught[pokemon_name[0]]
	if ok == false {
		return fmt.Errorf("You have not yet caught that pokemon")
	}
	fmt.Println("\nName:",p.Name)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("\nTypes:")
	for _, t:= range p.Types {
		fmt.Println(" - " + t.Type.Name)
	}

	currentLevel := GetLevelFromXP(p.CurrentXP)
	fmt.Printf("\nLevel: %d\n", currentLevel)
	fmt.Printf("Current XP: %d\n", p.CurrentXP)
	return nil
}
