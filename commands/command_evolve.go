package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

    "github.com/Ikit24/pokedexcli/internal/config"
)
func getJSONCached(cfg *config.Config, url string, out any) error {
	if data, ok := cfg.Cache.Get(url); ok {
		return json.Unmarshal(data, out)
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed: %d\n%s", res.StatusCode, string(body))
	}

	cfg.Cache.Add(url,body)
	return json.Unmarshal(body, out)
}

type SpeciesResponse struct {
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`
}

type EvolutionChain struct {
	Chain ChainLink `json:"chain"`
}
type ChainLink struct {
	Species struct {
		Name string `json:"name"`
	} `json:"species"`
    EvolvesTo []ChainLink `json:"evolves_to"`
    EvolutionDetails []struct {
        MinLevel *int `json:"min_level"`
    } `json:"evolution_details"`
}

}
func DetermineNextEvolution(cfg *config.Config, speciesName string) (string, int, error) {
	var sp SpeciesResponse
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", strings.ToLower(speciesName))
	if err := getJSONCached(cfg, url, &sp); err	!= nil {
		return "", 0, err
	}
	chainURL := sp.EvolutionChain.URL
	res, err := http.Get(url)
	if err != nil {
		return err
	}
}
