package main

import (
	"testing"
	"github.com/Ikit24/pokedexcli/internal/config"
)

func TestCleanInput(t *testing.T) {
	cases := [] struct {
		input    string
		expected []string
	}{
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "  hello  ",
			expected: []string{"hello"},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  HeLLo WorlD  ",
			expected: []string{"hello", "world"},
		},
	}


	for _,c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", c.input, actual, c.expected)
			}
		}
	}
}

func TestGetCommands(t *testing.T) {
	cfg := &config.Config{}
	commands := getCommands(cfg)

	expectedCommands := []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect", "pokedex", "battle", "save", "delete"}

	for _, cmdName := range expectedCommands {
		if _, exists := commands[cmdName]; !exists {
			t.Errorf("Expected command '%s' not found", cmdName)
		}
	}

	for cmdName, cmd := range commands {
		if cmd.Name == "" {
		t.Errorf("Command '%s' missing Name field", cmdName)
		}
		if cmd.Name == "" {
			t.Errorf("Command '%s' missing Description field", cmdName)
		}
		if cmd.Callback == nil {
			t.Errorf("Command '%s' missing Callback Field", cmdName)
		}
	}
}

func TestCleanInputEdgeCases(t *testing.T) {
	cases := []struct {
		input	 string
		expected []string
	}{
        {
            input:    "",
            expected: []string{},
        },
        {
            input:    "CATCH",
            expected: []string{"catch"},
        },
        {
            input:    "catch    pikachu",
            expected: []string{"catch", "pikachu"},
        },
        {
            input:    "\t\n  explore  area-name  \t\n",
            expected: []string{"explore", "area-name"},
        },
    }

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) = %v, expected %v", c.input, actual, c.expected)
				break
			}
		}
	}
}
