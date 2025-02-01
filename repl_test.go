package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		exp   []string
	}{
		{input: "", exp: nil},
		{input: "Hello world", exp: []string{"hello", "world"}},
		{input: "Pikachu mewtwo mew charizard", exp: []string{"pikachu", "mewtwo", "mew", "charizard"}},
		{input: "  ALL   UPPER CASE    ", exp: []string{"all", "upper", "case"}},
	}
	for _, c := range cases {
		res := cleanInput(c.input)
		if len(res) != len(c.exp) {
			t.Errorf("Input '%s' has %d results, expected %d", c.input, len(res), len(c.exp))
			t.Fatal()
		}
		for i := range res {
			if res[i] != c.exp[i] {
				t.Errorf("input: %s, exp: %s, got: %s", c.input, c.exp[i], res[i])
				t.Fatal()
			}
		}

	}
}
