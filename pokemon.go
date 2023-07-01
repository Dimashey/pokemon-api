package client

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) GetPokemonByName(ctx context.Context, name string) (Pokemon, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.apiURL+"/api/v2/pokemon/"+name,
		nil,
	)
	if err != nil {
		return Pokemon{}, PokemonFetchErr{Message: err.Error(), StatusCode: -1}
	}

	req.Header.Add("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, PokemonFetchErr{Message: err.Error(), StatusCode: -1}
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Pokemon{}, PokemonFetchErr{Message: "non-200 status code from the API", StatusCode: res.StatusCode}
	}

	var pokemon Pokemon

	err = json.NewDecoder(res.Body).Decode(&pokemon)

	if err != nil {
		return Pokemon{}, PokemonFetchErr{Message: err.Error(), StatusCode: res.StatusCode}
	}

	return pokemon, nil
}
