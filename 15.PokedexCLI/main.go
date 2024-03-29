package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
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
	Pokedex  map[string]pokeapi.Pokemon
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
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokémon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all the Pokémon caught",
			callback:    commandPokedex,
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
	const pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon/"

	if len(args) < 2 {
		return fmt.Errorf("catch command requires a valid Pokémon name")
	}
	fmt.Printf("Throwing a Pokéball at %s...\n", args[1])
	pokemonURL := pokemonBaseURL + args[1] + "/"
	pokemon, err := pokeapi.GetPokemon(pokemonURL, cf.Cache)
	if err != nil {
		return err
	}

	//assuming max base experience is 608 for a Pokémon currently
	catchChance := min(1, max(0, float64(pokemon.BaseExperience)/700.0))

	if rand.NormFloat64() > catchChance {
		fail_phrases := []string{
			"Oh no! The Pokémon broke free!",
			"Aww! It appeared to be caught!",
			"Shoot! It was so close too!",
			fmt.Sprintf("%s was caught! ... Just kidding!", pokemon.Name),
			"Almost had it!",
			fmt.Sprintf("%s escaped!", pokemon.Name),
		}
		fmt.Println(fail_phrases[rand.Intn(len(fail_phrases))])
		return nil
	}
	fmt.Printf("Caught %s!\n", pokemon.Name)
	//show the Pokémon caught
	for _, line := range pokemon.ASCIIArt {
		fmt.Println(line)
	}
	cf.Pokedex[pokemon.Name] = *pokemon
	return nil
}

func commandInspect(cf *appConfig, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf(("inspect command requires a valid Pokémon name"))
	}
	if pokemon, ok := cf.Pokedex[args[1]]; ok {

		pokemonStats := make([]string, 0)
		pokemonStats = append(pokemonStats, fmt.Sprintf("Name: %s", pokemon.Name))
		pokemonStats = append(pokemonStats, fmt.Sprintf("Height: %d", pokemon.Height))
		pokemonStats = append(pokemonStats, fmt.Sprintf("Weight: %d", pokemon.Weight))
		pokemonStats = append(pokemonStats, fmt.Sprintf("Base Experience: %d", pokemon.BaseExperience))
		pokemonStats = append(pokemonStats, "Stats:")
		for _, stat := range pokemon.Stats {
			pokemonStats = append(pokemonStats, fmt.Sprintf("  - %v: %d", stat.Stat.Name, stat.BaseStat))
		}
		pokemonStats = append(pokemonStats, "Types:")
		for _, t := range pokemon.Types {
			pokemonStats = append(pokemonStats, fmt.Sprintf("  - %s", t.Type.Name))
		}

		// fmt.Println("Name:", pokemon.Name)
		// fmt.Println("Height:", pokemon.Height)
		// fmt.Println("Weight:", pokemon.Weight)
		// fmt.Println("Stats:")
		// for _, stat := range pokemon.Stats {
		// 	fmt.Printf("  - %v: %d\n", stat.Stat.Name, stat.BaseStat)
		// }
		// fmt.Println("Types:")
		// for _, t := range pokemon.Types {
		// 	fmt.Println("  -", t.Type.Name)
		// }

		for i := 0; i < len(pokemonStats) || i < len(pokemon.ASCIIArt); i++ {
			stat := ""
			if i < len(pokemonStats) {
				stat = pokemonStats[i]
			}
			art := ""
			if i < len(pokemon.ASCIIArt) {
				art = pokemon.ASCIIArt[i]
			}
			fmt.Printf("%-30s %s\n", stat, art)
		}
	} else {
		fail_phrases := []string{
			"Who's that Pokémon?",
			"Could not find Pokémon in the Pokédex",
			"You need to catch the Pokémon first!",
			fmt.Sprintf("You have not caught %s yet!", args[1]),
		}
		fmt.Println(fail_phrases[rand.Intn(len(fail_phrases))])
	}
	return nil
}

func commandPokedex(cf *appConfig, _ []string) error {
	if len(cf.Pokedex) == 0 {
		fail_phrases := []string{
			"Pokedéx is empty! Go catch some Pokémon!",
			"No Pokémon caught yet!",
			"Your Pokédex is empty!",
			"Your Pokédex is empty! Go catch 'em all!",
		}
		fmt.Println(fail_phrases[rand.Intn(len(fail_phrases))])
		return nil
	}
	fmt.Println("Your Pokédex:")
	pokedex := make([]string, 0)
	for name := range cf.Pokedex {
		pokedex = append(pokedex, name)
	}
	for i := 0; i < len(pokedex); i += 2 {
		//2 columns
		for j := i; j < i+2 && j < len(pokedex); j++ {
			fmt.Printf(" - %-40s", pokedex[j])
		}
		fmt.Println()
		for k := 0; k < len(cf.Pokedex[pokedex[i]].ASCIIArt); k++ {
			for j := i; j < i+2 && j < len(pokedex); j++ {
				fmt.Printf("   %-43s", cf.Pokedex[pokedex[j]].ASCIIArt[k])
			}
			fmt.Println()
		}
		fmt.Println()

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
		Pokedex:  make(map[string]pokeapi.Pokemon),
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
		fmt.Print("\nPokedex > ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}

	cancel()
	//add a small delay to allow the go routines to exit gracefully
	time.Sleep(10 * time.Millisecond)
}
