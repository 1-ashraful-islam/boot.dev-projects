package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1-ashraful-islam/boot.dev-projects/15.PokedexCLI/internal/pokeapi"
	"github.com/1-ashraful-islam/boot.dev-projects/15.PokedexCLI/internal/pokecache"
)

type appConfig struct {
	Next     *string
	Previous *string
	Explore  *string
	Cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(cf *appConfig) error
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokémon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas in the Pokémon world.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area for all the Pokémon species",
			callback:    commandExplore,
		},
	}
}

func commandHelp(_ *appConfig) error {
	help_message := "\nWelcome to the Pokedex!\nUsages:\n"
	commands := getCliCommands()
	for _, command := range commands {
		help_message += fmt.Sprintf("  %s: %s\n", command.name, command.description)
	}
	fmt.Println(help_message)
	return nil
}

func commandExit(_ *appConfig) error {
	fmt.Println("Goodbye!")
	return nil
}

func commandMap(cf *appConfig) error {
	if cf.Next == nil {
		return fmt.Errorf("no more location areas")
	}
	locations, err := pokeapi.GetLocationArea(*cf.Next, cf.Cache)
	if err != nil {
		return err
	}

	cf.Next = locations.Next
	cf.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(cf *appConfig) error {
	if cf.Previous == nil {
		return fmt.Errorf("no previous location areas")
	}
	locations, err := pokeapi.GetLocationArea(*cf.Previous, cf.Cache)
	if err != nil {
		return err
	}

	cf.Next = locations.Next
	cf.Previous = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cf *appConfig) error {
	const exploreBaseURL = "https://pokeapi.co/api/v2/location-area/"
	//reset the explore config value to nil after the command is executed
	defer func() { cf.Explore = nil }()

	if cf.Explore == nil {
		return fmt.Errorf("explore command requires a valid location area name")
	}

	fmt.Printf("Exploring %v ...\n", *cf.Explore)

	exploreURL := exploreBaseURL + *cf.Explore
	encounters, err := pokeapi.GetPokemonInLocationArea(exploreURL, cf.Cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokémon:")
	for _, encounter := range encounters.PokemonEncounters {
		fmt.Println(" - ", encounter.Pokemon.Name)
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	locationURL := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	commands := getCliCommands()
	config := &appConfig{
		Next:     &locationURL,
		Previous: nil,
		Cache:    pokecache.NewCache(ctx, 360*time.Second),
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		command_str := strings.Split(line, " ")[0]
		command, ok := commands[command_str]
		if ok {
			if command.name == "explore" && len(strings.Split(line, " ")) > 1 {
				config.Explore = &strings.Split(line, " ")[1]
			}

			err := command.callback(config)
			if err != nil {
				fmt.Println("Error:", err)
			}
			if command.name == "exit" {
				break
			}
		} else {
			fmt.Printf("Unknown command: %s\n", command_str)
			commands["help"].callback(nil)
		}
		fmt.Print("Pokedex > ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}

	cancel()
	//add a small delay to allow the go routines to exit
	time.Sleep(10 * time.Millisecond)
}
