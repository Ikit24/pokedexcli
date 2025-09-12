package pokeapi

import (
    "time"
)

type LocationAreasResponse struct {
	Next     string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Count    int    `json:"count"`
}

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type BattlePokemon struct {
    Name      string   `json:"name"`
    Types     []struct {
        Type struct {
            Name string `json:"name"`
        } `json:"type"`
    } `json:"types"`
    Stats     []struct {
        BaseStat int `json:"base_stat"`
        Stat     struct {
            Name string `json:"name"`
        } `json:"stat"`
    } `json:"stats"`
    Abilities []struct {
        Ability struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"ability"`
        IsHidden bool `json:"is_hidden"`
    } `json:"abilities"`
    Moves []struct {
        Move struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"move"`
    } `json:"moves"`
    StatusEffects []string `json:"status_effects"`
    EvolvesTo string `json:"evolves_to"`
    CaughtAt time.Time `json:"caught_at"`
    ID                   int `json:"id"`
    BaseExperience       int `json:"base_experience"`
    Height               int `json:"height"`
    Weight               int `json:"weight"`
    CurrentHP            int `json:"current_hp"`
    CurrentXP            int `json:"current_xp"`
    Level                int `json:"level"`
    EvolutionDelaySecs   int `json:"evolution_delay_secs"`
    MinLevelForEvolution int `json:"min_level_for_evolution"`
    HasEvolved bool `json:"has_evolved"`
}

type EvolutionChain struct {
Chain ChainLink `json:"chain"`
}

type ChainLink struct {
    EvolvesTo []struct {
        Species struct {
            Name string `json:"name"`
        } `json:"species"`
        EvolutionDetails []struct {
            MinLevel *int `json:"min_level"`
        } `json:"evolution_details"`
        EvolvesTo []ChainLink `json:"evolves_to"`
    } `json:"evolves_to"`
    Species struct {
        Name string `json:"name"`
    } `json:"species"`
}
