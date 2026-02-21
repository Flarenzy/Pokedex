package cmd

import (
	"strings"
	"testing"
)

func TestNewCommands(t *testing.T) {
	commands := NewCommands()
	tests := []struct {
		name string
	}{
		{name: "help"},
		{name: "exit"},
		{name: "map"},
		{name: "mapb"},
		{name: "explore"},
		{name: "catch"},
		{name: "inspect"},
		{name: "pokedex"},
	}

	if len(commands) != len(tests) {
		t.Fatalf("expected %d commands, got %d", len(tests), len(commands))
	}

	for _, tc := range tests {
		command, ok := commands[tc.name]
		if !ok {
			t.Fatalf("missing command %q", tc.name)
		}
		if command == nil {
			t.Fatalf("command %q is nil", tc.name)
		}
		if command.description == "" {
			t.Fatalf("description for %q is empty", tc.name)
		}
		if command.Callback == nil {
			t.Fatalf("callback for %q is nil", tc.name)
		}
	}
}

func TestNewCommandsHelpTextStable(t *testing.T) {
	_ = NewCommands()
	first := helpText
	_ = NewCommands()
	second := helpText

	if second != first {
		t.Fatalf("expected helpText to stay stable across calls")
	}

	tests := []struct {
		name string
	}{
		{name: "help"},
		{name: "exit"},
		{name: "map"},
		{name: "mapb"},
		{name: "explore"},
		{name: "catch"},
		{name: "inspect"},
		{name: "pokedex"},
	}

	for _, tc := range tests {
		prefix := tc.name + ": "
		count := strings.Count(helpText, prefix)
		if count != 1 {
			t.Fatalf("expected %q to appear exactly once in helpText, got %d", prefix, count)
		}
	}
}
