package client

import (
	"errors"
	"fmt"
)

var ErrFetchingPokemon = errors.New("failed to fetch")

type PokemonFetchErr struct {
	StatusCode int
	Message    string
}

func (e PokemonFetchErr) Error() string {
	return fmt.Sprintf("failed to fetch pokemon: %s with status code %d", e.Message, e.StatusCode)
}
