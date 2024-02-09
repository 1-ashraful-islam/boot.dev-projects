package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/1-ashraful-islam/boot.dev-projects/15.PokedexCLI/internal/pokeapi"
	"github.com/1-ashraful-islam/boot.dev-projects/15.PokedexCLI/internal/pokecache"

	tea "github.com/charmbracelet/bubbletea"
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
	callback    func(cf *appConfig, args []string) (string, error)
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

func commandHelp(_ *appConfig, _ []string) (string, error) {
	help_message := "\nWelcome to the Pokedex!\nUsages:\n"
	commands := getCliCommands()
	for _, command := range commands {
		help_message += fmt.Sprintf("  %s: %s\n", command.name, command.description)
	}
	return help_message, nil
}

func commandExit(_ *appConfig, _ []string) (string, error) {
	return "Goodbye!", nil
}

func commandMap(cf *appConfig, _ []string) (string, error) {
	if cf.Next == nil {
		return "", fmt.Errorf("no more location areas")
	}
	locations, err := pokeapi.GetLocationArea(*cf.Next, cf.Cache)
	if err != nil {
		return "", err
	}

	cf.Next = locations.Next
	cf.Previous = locations.Previous

	result := ""
	for _, location := range locations.Results {
		result += location.Name + "\n"
	}
	return result, nil
}

func commandMapb(cf *appConfig, _ []string) (string, error) {
	if cf.Previous == nil {
		return "", fmt.Errorf("no previous location areas")
	}
	locations, err := pokeapi.GetLocationArea(*cf.Previous, cf.Cache)
	if err != nil {
		return "", err
	}

	cf.Next = locations.Next
	cf.Previous = locations.Previous

	result := ""
	for _, location := range locations.Results {
		result += location.Name + "\n"
	}

	return result, nil
}

func commandExplore(cf *appConfig, args []string) (string, error) {
	const exploreBaseURL = "https://pokeapi.co/api/v2/location-area/"

	if len(args) < 2 {
		return "", fmt.Errorf("explore command requires a valid location area name")
	}

	result := fmt.Sprintf("Exploring %v ...\n", args[1])

	exploreURL := exploreBaseURL + args[1] + "/"
	encounters, err := pokeapi.GetPokemonInLocationArea(exploreURL, cf.Cache)
	if err != nil {
		return "", err
	}

	result += "Found Pokémon:\n"
	for _, encounter := range encounters.PokemonEncounters {
		result += " - " + encounter.Pokemon.Name + "\n"
	}

	return result, nil
}

func commandCatch(cf *appConfig, args []string) (string, error) {
	const pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon/"

	if len(args) < 2 {
		return "", fmt.Errorf("catch command requires a valid Pokémon name")
	}
	result := fmt.Sprintf("Throwing a Pokéball at %s...\n", args[1])
	pokemonURL := pokemonBaseURL + args[1] + "/"
	pokemon, err := pokeapi.GetPokemon(pokemonURL, cf.Cache)
	if err != nil {
		return "", err
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
		result += fail_phrases[rand.Intn(len(fail_phrases))]
		return result, nil
	}
	result += fmt.Sprintf("Caught %s!\n", pokemon.Name)
	//show the Pokémon caught
	for _, line := range pokemon.ASCIIArt {
		result += line + "\n"
	}
	cf.Pokedex[pokemon.Name] = *pokemon
	return result, nil
}

func commandInspect(cf *appConfig, args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("inspect command requires a valid Pokémon name")
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

		result := ""
		for i := 0; i < len(pokemonStats) || i < len(pokemon.ASCIIArt); i++ {
			stat := ""
			if i < len(pokemonStats) {
				stat = pokemonStats[i]
			}
			art := ""
			if i < len(pokemon.ASCIIArt) {
				art = pokemon.ASCIIArt[i]
			}
			result += fmt.Sprintf("%-30s %s\n", stat, art)
		}
		return result, nil
	} else {
		fail_phrases := []string{
			"Who's that Pokémon?",
			"Could not find Pokémon in the Pokédex",
			"You need to catch the Pokémon first!",
			fmt.Sprintf("You have not caught %s yet!", args[1]),
		}
		return "", fmt.Errorf(fail_phrases[rand.Intn(len(fail_phrases))])
	}
}

func commandPokedex(cf *appConfig, _ []string) (string, error) {
	if len(cf.Pokedex) == 0 {
		fail_phrases := []string{
			"Pokedéx is empty! Go catch some Pokémon!",
			"No Pokémon caught yet!",
			"Your Pokédex is empty!",
			"Your Pokédex is empty! Go catch 'em all!",
		}
		return "", fmt.Errorf(fail_phrases[rand.Intn(len(fail_phrases))])
	}
	result := "Your Pokédex:\n"
	pokedex := make([]string, 0)
	for name := range cf.Pokedex {
		pokedex = append(pokedex, name)
	}
	for i := 0; i < len(pokedex); i += 2 {
		//2 columns
		for j := i; j < i+2 && j < len(pokedex); j++ {
			result += fmt.Sprintf(" - %-40s", pokedex[j])
		}
		result += "\n"
		for k := 0; k < len(cf.Pokedex[pokedex[i]].ASCIIArt); k++ {
			for j := i; j < i+2 && j < len(pokedex); j++ {
				result += fmt.Sprintf("   %-43s", cf.Pokedex[pokedex[j]].ASCIIArt[k])
			}
			result += "\n"
		}
		result += "\n"

	}
	return result, nil
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

	startingScreen := &MainScreen{
		config:      config,
		cliCommands: commands,
		cmds:        []string{"map", "mapb", "explore", "catch", "inspect", "pokedex", "exit"},
		cursor:      0,
	}
	app := tea.NewProgram(&appModel{
		config:        config,
		cliCommands:   commands,
		screens:       make(map[string]UIScreen),
		currentScreen: startingScreen,
	})

	if _, err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cancel()
	//add a small delay to allow the go routines to exit gracefully
	time.Sleep(10 * time.Millisecond)
}

// interface to separate UI screen or context
type UIScreen interface {
	Update(msg tea.Msg) (UIScreen, tea.Cmd)
	View() string
}

type MainScreen struct {
	config      *appConfig
	cliCommands map[string]cliCommand
	cmds        []string
	cursor      int
}

func (m *MainScreen) Update(msg tea.Msg) (UIScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
			// left keys
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		// right keys
		case "right":
			if m.cursor < len(m.cmds)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m *MainScreen) View() string {
	var s string

	cmdHelp := ""
	//main menu
	for i, command := range m.cmds {
		cursorBegin := " "
		cursorEnd := " "

		if m.cursor == i {
			cursorBegin = "<"
			cursorEnd = ">"
			cmdHelp = m.cliCommands[command].description
		}
		s += fmt.Sprintf(" %s%s%s ", cursorBegin, command, cursorEnd)
	}

	s += "\n\n" + cmdHelp

	//footer
	s += "\nPress 'q' or 'esc' to quit.\n"
	return "Welcome to Pokédex!\n" + s
}

type appModel struct {
	config        *appConfig
	cliCommands   map[string]cliCommand
	screens       map[string]UIScreen
	currentScreen UIScreen
	prevScreen    UIScreen
}

// Init implements tea.Model.
func (m *appModel) Init() tea.Cmd {
	return nil
}

// delegate the update to the current screen
func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.currentScreen, cmd = m.currentScreen.Update(msg)
	return m, cmd
}

// delegate the view to the current screen
func (m *appModel) View() string {
	return m.currentScreen.View()
}
