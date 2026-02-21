package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mustReadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "pokeapi", name))
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return b
}

func TestMainE2EWithFixtures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e in short mode")
	}

	mapPage := mustReadFixture(t, "location-area-page1.json")
	areaOne := mustReadFixture(t, "location-area-1.json")
	pokemonMew := mustReadFixture(t, "pokemon-mew.json")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/location-area/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(mapPage)
		case "/api/v2/location-area/1/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(areaOne)
		case "/api/v2/pokemon/mew":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(pokemonMew)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	script := strings.Join([]string{
		"help",
		"map",
		"explore 1/",
		"catch mew",
		"exit",
	}, "\n") + "\n"

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"POKEDEX_AREA_URL="+ts.URL+"/api/v2/location-area/",
		"POKEDEX_POKEMON_URL="+ts.URL+"/api/v2/pokemon/",
	)
	cmd.Stdin = strings.NewReader(script)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		t.Fatalf("app run failed: %v\noutput:\n%s", err, output.String())
	}

	out := output.String()
	expectedContains := []string{
		"Welcome to the Pokedex!",
		"canalave-city-area",
		"Exploring area:  1/",
		"Pokemon #1:",
		"Throwing a Pokeball at mew...",
		"Closing the Pokedex... Goodbye!",
	}
	for _, s := range expectedContains {
		if !strings.Contains(out, s) {
			t.Fatalf("expected output to contain %q\nfull output:\n%s", s, out)
		}
	}
}
