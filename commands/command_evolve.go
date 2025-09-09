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

	cfg.Cache.Add(url, body)
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

func DetermineNextEvolution(cfg *config.Config, speciesName string) (string, int, error) {
	var sp SpeciesResponse
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", strings.ToLower(speciesName))
	if err := getJSONCached(cfg, url, &sp); err != nil {
		return "", 0, err
	}
	
	var chain EvolutionChain
	if err := getJSONCached(cfg, sp.EvolutionChain.URL, &chain); err != nil {
		return "", 0, err
	}
	
	node := findNode(chain.Chain, strings.ToLower(speciesName))
	if node == nil || len(node.EvolvesTo) == 0 {
		return "", 0, nil
	}
	
	child := node.EvolvesTo[0]
	next := child.Species.Name
	min := 1
	for _, d := range child.EvolutionDetails {
		if d.MinLevel != nil {
			min = *d.MinLevel
			break
		}
	}
	return next, min, nil
}

func findNode(n ChainLink, target string) *ChainLink {
	if n.Species.Name == target {
		return &n
	}
	for i := range n.EvolvesTo {
		if res := find.Node(n.EvolvesTo[i], target); res != nil {
			return res
		}
	}
	return nil
}

func ReadyToEvolve(p pokeapi.BattlePokemon, now time.Time) bool {
    if p.HasEvolved || p.EvolvesTo == "" { return false }
    if p.Level < p.MinLevelForEvolution { return false }
    if now.Sub(p.CaughtAt) < time.Duration(p.EvolutionDelaySecs)*time.Second { return false }
    return true
}
