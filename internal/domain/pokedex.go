package domain

import "github.com/Flarenzy/Pokedex/internal/pokedex"

type Pokedexer interface {
	Add(p pokedex.Pokemon)
	Remove(p pokedex.Pokemon)
	GetAllPokemon() []pokedex.Pokemon
	GetPokemonByName(name string) (pokedex.Pokemon, error)
}
