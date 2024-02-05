package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/1-ashraful-islam/boot.dev-projects/15.PokedexCLI/internal/pokecache"
)

type locationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemonInArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationArea(url string, cf *pokecache.Cache) (*locationArea, error) {

	body, err := getPokemonData(url, cf)
	if err != nil {
		return nil, err
	}

	locations := &locationArea{}
	if err := json.Unmarshal(body, locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func GetPokemonInLocationArea(url string, cf *pokecache.Cache) (*pokemonInArea, error) {
	body, err := getPokemonData(url, cf)
	if err != nil {
		return nil, err
	}

	pokemons := &pokemonInArea{}
	if err := json.Unmarshal(body, pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

func getPokemonData(url string, cf *pokecache.Cache) ([]byte, error) {
	if val, ok := cf.Get(url); ok {
		return val, nil
	}

	body, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	cf.Add(url, body)
	return body, nil
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
