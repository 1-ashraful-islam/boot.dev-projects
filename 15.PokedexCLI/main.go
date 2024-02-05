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
	callback    func(cf *appConfig, args []string) error
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
		"catch": {
			name:        "catch",
			description: "Catch a Pokémon",
			callback:    commandCatch,
		},
	}
}

func commandHelp(_ *appConfig, _ []string) error {
	help_message := "\nWelcome to the Pokedex!\nUsages:\n"
	commands := getCliCommands()
	for _, command := range commands {
		help_message += fmt.Sprintf("  %s: %s\n", command.name, command.description)
	}
	fmt.Println(help_message)
	return nil
}

func commandExit(_ *appConfig, _ []string) error {
	fmt.Println("Goodbye!")
	return nil
}

func commandMap(cf *appConfig, _ []string) error {
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

func commandMapb(cf *appConfig, _ []string) error {
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

func commandExplore(cf *appConfig, args []string) error {
	const exploreBaseURL = "https://pokeapi.co/api/v2/location-area/"

	if len(args) < 2 {
		return fmt.Errorf("explore command requires a valid location area name")
	}

	fmt.Printf("Exploring %v ...\n", args[1])

	exploreURL := exploreBaseURL + args[1] + "/"
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

func commandCatch(cf *appConfig, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("catch command requires a valid Pokémon name")
	}
	fmt.Printf("Throwing a Pokéball at %s...\n", args[1])
	return fmt.Errorf("not implemented")
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
		args := strings.Split(scanner.Text(), " ")
		command, ok := commands[args[0]]
		if ok {
			if err := command.callback(config, args); err != nil {
				fmt.Println("Error:", err)
			}
			if command.name == "exit" {
				break
			}
		} else {
			fmt.Printf("Unknown command: %s\n", args[0])
			commands["help"].callback(nil, nil)
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
