package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientCanHitAPI(t *testing.T) {
	t.Run("happy path - can hit the api and return a pokemon", func(t *testing.T) {
		myClient := New()

		poke, err := myClient.GetPokemonByName(context.Background(), "pikachu")

		assert.NoError(t, err)
		assert.Equal(t, "pikachu", poke.Name)
	})

	t.Run("sad path - return an error when the pokemon does not exists", func(t *testing.T) {
		myClient := New()

		_, err := myClient.GetPokemonByName(context.Background(), "non-existant-pokemon")

		assert.Error(t, err)
		assert.Equal(t, PokemonFetchErr{Message: "non-200 status code from the API", StatusCode: 404}, err)
	})

	t.Run("happy path - testing the WithAPIURL option function", func(t *testing.T) {
		myClient := New(WithAPIURL("my-test-url"))

		assert.Equal(t, "my-test-url", myClient.apiURL)
	})

	t.Run("happy path - tests with httpclient works", func(t *testing.T) {
		myClient := New(WithAPIURL("my-test-url"), WithHTTPClient(&http.Client{Timeout: 1 * time.Second}))

		assert.Equal(t, "my-test-url", myClient.apiURL)
		assert.Equal(t, 1*time.Second, myClient.httpClient.Timeout)
	})

	t.Run("happy test - able to hit locally running test server", func(t *testing.T) {
		ts := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					fmt.Fprintf(w, `{"name": "pikachu", "height": 10}`)
				},
			),
		)

		defer ts.Close()

		myClient := New(WithAPIURL(ts.URL))

		poke, err := myClient.GetPokemonByName(context.Background(), "pikachu")

		assert.NoError(t, err)
		assert.Equal(t, 10, poke.Height)
	})

	t.Run("sad path test - able to handle 500 status from the API", func(t *testing.T) {
		ts := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				},
			),
		)

		defer ts.Close()

		myClient := New(WithAPIURL(ts.URL))

		_, err := myClient.GetPokemonByName(context.Background(), "pikachu")

		assert.Error(t, err)
	})
}
