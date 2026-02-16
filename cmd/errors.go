package cmd

import (
	"errors"
)

var ErrStop = errors.New("stop")
var ErrNoPokemon = errors.New("no pokemon")
