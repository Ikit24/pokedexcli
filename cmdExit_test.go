package main

import (
	"testing"
	"fmt"
    "github.com/Ikit24/pokedexcli/internal/config"
    "github.com/Ikit24/pokedexcli/commands"
)

func TestExit (t *testing.T) {
	cases := [] struct {
		input	 []string
		expected error
	}{
		{
			input: []string{"exit", "and", "save"},
			expected: fmt.Errorf("exit command doesn't require any arguments."),
		},
	}

	cfg := &config.Config{}

	for _, c := range cases {
		actual := commands.CommandExit(cfg, c.input)
		if actual.Error() != c.expected.Error() {
			t.Errorf("Test failed, CommandExit returned: '%s', expected: '%s'", actual, c.expected)
		}
	}
}
