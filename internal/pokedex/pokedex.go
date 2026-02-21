package pokedex

import (
	"errors"
	"sync"
)

var ErrPokemonNotFound = errors.New("pokemon not found")

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type Pokedex struct {
	Owned map[string]Pokemon
	mu    sync.RWMutex
}

func (p *Pokedex) Add(pokemon Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Owned[pokemon.Name] = pokemon
}

func (p *Pokedex) Remove(pokemon Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.Owned[pokemon.Name]
	if ok {
		delete(p.Owned, pokemon.Name)

	}
}
func (p *Pokedex) GetAllPokemon() []Pokemon {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allPokemon := make([]Pokemon, 0)
	for _, pokemon := range p.Owned {
		allPokemon = append(allPokemon, pokemon)
	}
	return allPokemon
}

func (p *Pokedex) GetPokemonByName(name string) (Pokemon, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, pokemon := range p.Owned {
		if pokemon.Name == name {
			return pokemon, nil
		}
	}
	return Pokemon{}, ErrPokemonNotFound
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		Owned: make(map[string]Pokemon),
		mu:    sync.RWMutex{},
	}

}
