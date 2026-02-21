package pokedex

import (
	"errors"
	"testing"
)

func TestPokedexCRUD(t *testing.T) {
	t.Parallel()

	p := NewPokedex()
	mew := Pokemon{Name: "mew", Height: 4, Weight: 40}
	pikachu := Pokemon{Name: "pikachu", Height: 4, Weight: 60}

	p.Add(mew)
	p.Add(pikachu)

	got, err := p.GetPokemonByName("mew")
	if err != nil {
		t.Fatalf("expected mew, got error: %v", err)
	}
	if got.Name != "mew" {
		t.Fatalf("expected mew, got %q", got.Name)
	}

	all := p.GetAllPokemon()
	if len(all) != 2 {
		t.Fatalf("expected 2 pokemon, got %d", len(all))
	}

	p.Remove(mew)
	_, err = p.GetPokemonByName("mew")
	if !errors.Is(err, ErrPokemonNotFound) {
		t.Fatalf("expected ErrPokemonNotFound, got %v", err)
	}

	p.Remove(Pokemon{Name: "does-not-exist"})
	all = p.GetAllPokemon()
	if len(all) != 1 {
		t.Fatalf("expected 1 pokemon after remove, got %d", len(all))
	}
}

func TestGetPokemonByNameNotFound(t *testing.T) {
	t.Parallel()

	p := NewPokedex()
	_, err := p.GetPokemonByName("missing")
	if !errors.Is(err, ErrPokemonNotFound) {
		t.Fatalf("expected ErrPokemonNotFound, got %v", err)
	}
}
