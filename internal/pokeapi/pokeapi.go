package pokeapi

type LocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type BattlePokemon struct {
    Name           string `json:"name"`
    ID             int    `json:"id"`
    BaseExperience int    `json:"base_experience"`
    Height         int    `json:"height"`
    Weight         int    `json:"weight"`
 
    Stats []struct {
        BaseStat int `json:"base_stat"`
        Stat     struct {
            Name string `json:"name"`
        } `json:"stat"`
    } `json:"stats"`
 
    Types []struct {
        Type struct {
            Name string `json:"name"`
        } `json:"type"`
    } `json:"types"`
 
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
 
    CurrentHP    int      `json:"-"`
    StatusEffects []string `json:"-"`
    CurrentXP    int      `json:"-"`
    Level        int      `json:"-"`
}
