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

func DetermineNextEvolution(ctx context.Context, client *http.Client, speciesName string) (string, int, error) {
	pass
}
