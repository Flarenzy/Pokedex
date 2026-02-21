package cmd

import "errors"

var (
	ErrStop               = errors.New("stop")
	ErrNoPokemon          = errors.New("no pokemon")
	ErrNoPokemonToInspect = errors.New("no pokemon to inspect")
	ErrNoAreaToExplore    = errors.New("no command to explore")
	ErrEmptyPokedex       = errors.New("no pokedex found")
)
