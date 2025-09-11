package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
    "github.com/Ikit24/pokedexcli/internal/config"
	"github.com/Ikit24/pokedexcli/internal/pokeapi"
)
func getJSONCached(cfg *config.Config, url string, out any) error {
	if data, ok := cfg.Cache.Get(url); ok {
		return json.Unmarshal(data, out)
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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
	if strings.ToLower(n.Species.Name) == strings.ToLower(target) {
		return &n
	}
	for i := range n.EvolvesTo {
		if res := findNode(n.EvolvesTo[i], target); res != nil {
			return res
		}
	}
	return nil
}

func ReadyToEvolve(p pokeapi.BattlePokemon, now time.Time)   bool {
    if p.HasEvolved || p.EvolvesTo == "" { return false }
    if p.Level < p.MinLevelForEvolution { return false }
    if now.Sub(p.CaughtAt) < time.Duration(p.EvolutionDelaySecs)*time.Second { return false }
    return true
}

func EvolveTo(cfg *config.Config, p *pokeapi.BattlePokemon, next string) error {
	if next == "" {
		return nil
	}
	var evolved pokeapi.BattlePokemon
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", strings.ToLower(next))
	if err := getJSONCached(cfg, pokemonURL, &evolved); err != nil {
		return err
	}

	lvl, xp := p.Level, p.CurrentXP
	status, currHP := p.StatusEffects, p.CurrentHP

	*p = evolved
	p.Level = lvl
	p.CurrentXP = xp
	p.StatusEffects = status

	maxHP, _ := getStatValue(p.Stats, "hp")
    if currHP > maxHP {
		currHP = maxHP
	}
    p.CurrentHP = currHP

	p.HasEvolved = true
	p.EvolvesTo = ""
	p.MinLevelForEvolution = 0
	return nil
}

func RunEvolutionPass(cfg *config.Config) ([]string, error) {
	now := time.Now().UTC()
	msgs := []string{}
	toInsert := map[string]pokeapi.BattlePokemon{}
	toDelete := []string{}

	for key, p := range cfg.Caught {
		if p.EvolvesTo == "" && !p.HasEvolved {
			if next, min, err := DetermineNextEvolution(cfg, p.Name); err == nil {
				p.EvolvesTo, p.MinLevelForEvolution = next, min
			}
		}
		if ReadyToEvolve(p, now) {
			oldName := p.Name
			if err := EvolveTo(cfg, &p, p.EvolvesTo); err != nil {
				return msgs, err
			}
			newKey := p.Name
			toDelete = append(toDelete, key)
			toInsert[newKey] = p
			msgs = append(msgs, fmt.Sprintf("%s -> %s", oldName, p.Name))
		} else {
			cfg.Caught[key] = p
		}
	}
	for _, k := range toDelete { delete(cfg.Caught, k) }
	for k, v := range toInsert { cfg.Caught[k] = v }
	return msgs, nil
}

func CommandEvolve(cfg *config.Config, args []string) error {
	msgs, err := RunEvolutionPass(cfg)
	if err != nil {
		return  fmt.Errorf("evolution check failed: %w", err)
	}
	if len(msgs) == 0 {
		fmt.Println("No Pokemon ready to evolve.")
		return nil
	}
	for _, m := range msgs {
		fmt.Println(m)
	}
	AutoSave(cfg)
	return nil
}
