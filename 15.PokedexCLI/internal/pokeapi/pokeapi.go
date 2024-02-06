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

type Pokemon struct {
	Name                   string `json:"name"`
	BaseExperience         int    `json:"base_experience"`
	Height                 int    `json:"height"`
	Weight                 int    `json:"weight"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Sprites                struct {
		Others struct {
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
	} `json:"sprites"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
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

func GetPokemon(url string, cf *pokecache.Cache) (*Pokemon, error) {
	body, err := getPokemonData(url, cf)
	if err != nil {
		return nil, err
	}

	pokemon := &Pokemon{}
	if err := json.Unmarshal(body, pokemon); err != nil {
		return nil, err
	}

	return pokemon, nil
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
