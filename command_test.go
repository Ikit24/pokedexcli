package main

import (
	"testing"
	"time"
    "github.com/Ikit24/pokedexcli/internal/config"
    "github.com/Ikit24/pokedexcli/commands"
	"github.com/Ikit24/pokedexcli/internal/pokecache"
)

func TestExit (t *testing.T) {
	cases := [] struct {
		input	 []string
		expected string
	}{
		{
			input: []string{"exit", "and", "save"},
			expected: "exit command doesn't require any arguments.",
		},
	}

	cfg := &config.Config{}

	for _, c := range cases {
		actual := commands.CommandExit(cfg, c.input)
		if actual.Error() != c.expected {
			t.Errorf("Test failed, CommandExit returned: '%s', expected: '%s'", actual.Error(), c.expected)
		}
	}
}

func TestExplore (t *testing.T) {
	cases := [] struct {
		next	 string
		input	 []string
		expected string
	}{
		{
			next:	"",
			input: []string{"explore"},
			expected: "You must map the area first! Use the 'map' command.",
		},
		{
			next:	"https://pokeapi.co/api/v2/location-area/",
			input: []string{"explore"},
			expected: "You must map the area first! Use the 'map' command.",
		},
	}

	for _, c := range cases {
		cfg := &config.Config{
			Next: c.next,
		}
		actual := commands.CommandExplore(cfg, c.input)
		if actual.Error() != c.expected {
			t.Errorf("Test failed, CommandExplore returned: '%s', expected: '%s'", actual.Error(), c.expected)
		}
	}
}

 func TestExploreAreaName(t *testing.T) {
	cases := [] struct {
		input	 []string
		expected string
	}{
		{
			input: []string{},
			expected: "Must provide location area name!",
		},
	}

	for _, c := range cases {
		cfg := &config.Config{
			Next: "https://pokeapi.co/api/v2/location-area/some-valid-area",
			Cache: pokecache.NewCache(1 * time.Minute),
		}
		actual := commands.CommandExplore(cfg, c.input)
		if actual.Error() != c.expected {
			t.Errorf("Test failed, CommandExplore returned: '%s', expected: '%s'", actual.Error(), c.expected)
		}
	}
}

func TestInspect(t *testing.T) {
	cases := [] struct {
		input	 []string
		expected string
	}{
		{
			input: []string{},
			expected: "Must provide pokemon name in order to display it's details",
		},
	}

	for _, c :=  range cases {
		cfg := &config.Config{}
		actual := commands.CommandInspect(cfg, c.input)
		if actual.Error() != c.expected {
			t.Errorf("Test failed, CommandInspect returned: '%s', expected: '%s'", actual.Error(), c.expected)
		}
	}
}

func TestBattle(t *testing.T) {
	cases := [] struct {
		input	 []string
		expected string
	}{
		{
			input: []string{},
			expected: "Usage: battle <your-pokemon> <target-pokemon>",
		},
	}
	for _, c :=  range cases {
		cfg := &config.Config{}
		actual := commands.CommandBattle(cfg, c.input)
		if actual.Error() != c.expected {
			t.Errorf("Test failed, CommandBattle returned: '%s', expected: '%s'", actual.Error(), c.expected)
		}
	}
}
