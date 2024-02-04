package pokeapi

import (
	"encoding/json"
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

func GetLocationArea(url string, cf *pokecache.Cache) (*locationArea, error) {

	var body []byte

	if val, ok := cf.Get(url); ok {
		body = val
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		cf.Add(url, body)
	}

	locations := &locationArea{}
	err := json.Unmarshal(body, locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
