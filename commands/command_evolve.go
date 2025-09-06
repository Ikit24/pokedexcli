package commands

func getJSON(ctx context.Context, client *http.Client, url string, any) error {
	req, _ := http.NewRequestWirthContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {return err}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func DetermineNextEvolution(ctx context.Context, client *http.Client, speciesName string) (string, int, error) {
	pass
}
