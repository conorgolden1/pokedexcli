package main

import (
	"bufio"
	"errors"
	"fmt"
	"internal/pokecache"
	"os"
	"reflect"
)

func getCliCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name: "map",
			description: `Explore the pokemon world!
            Displays the names of 20 location areas in the Pokemon World
            Each subsqequent call to map will display the next 20 locations.`,
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: `Map Back: Explore the previous pokemon world
            Each subsequent call to map will display the previous 20 locations

            This will error if on at map cursor is at the beginning`,
			callback: commandMapb,
		},
		"explore": {
			name: "explore",
			description: `explore <name>
            List all of the Pokemon in a given area`,
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: `catch <name>
            Catch a pokemon and add it to the your pokedex`,
			callback: commandCatch,
		},
        "inspect": {
			name: "inspect",
			description: `inspect <name>
            Inspect a pokemon that you have caught`,
			callback: commandInspect,
		},
        "pokedex": {
			name: "pokedex",
			description: `pokedex <name>
            Display all of the pokemon you have caught so far`,
            callback: commandPokedex,
		},
	}
}

func getCommand(scanner *bufio.Scanner) string {
	fmt.Printf(" pokedex > ")
	scanner.Scan()
	fmt.Printf("\n")
	text := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input from stdin")
		os.Exit(1)
	}
	return text
}

type Config struct {
	args      []string
	cache     pokecache.Cache
	mapObject Map
	pokedex   map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func commandHelp(_ *Config) error {
	fmt.Printf(`Welcome to the Pokedex!
    Usage:

    help: Displays a help message
    exit: Exit the Pokedex`)
	fmt.Println()
	return nil
}

func commandExit(_ *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(c *Config) error {
	if (reflect.DeepEqual(c.mapObject, Map{})) {
		mapObject, err := getMap("", &c.cache)
		if err != nil {
			return err
		}
		c.mapObject = mapObject
	} else {
		next, err := c.mapObject.nextMap(&c.cache)
		if err != nil {
			return err
		}
		c.mapObject = next
	}
	c.mapObject.Print()
	return nil
}

func commandMapb(c *Config) error {
	if (reflect.DeepEqual(c.mapObject, Map{}) || c.mapObject.Previous == nil) {
		return errors.New("\tUnable to display previous map: At the beginning of the map stream\n")
	}
	prev, err := c.mapObject.prevMap(&c.cache)

	if err != nil {
		return err
	}

	c.mapObject = prev
	c.mapObject.Print()

	return nil
}

func commandExplore(c *Config) error {
	if len(c.args) == 0 {
		return errors.New("Not enough arguements\nUsage: explore <location-name1> <location-name2>...")
	}

	for _, v := range c.args {
		fmt.Printf("Exploring %v...\n", v)
		loc, err := GetLocation(v, &c.cache)
		if err != nil {
			fmt.Println(err)
			continue
		}
		loc.Print()
	}
	return nil
}

func commandCatch(c *Config) error {
	if len(c.args) == 0 {
		return errors.New("Not enough arguements\nUsage: catch <pokemon-name>")
	}

	if len(c.args) > 1 {
		return errors.New("Too many arguements\nUsage: catch <pokemon-name>")
	}

	pm, err := GetPokemon(c.args[0], &c.cache)
	if err != nil {
		return err
	}

	pm.Catch(c)

	return nil
}

func commandInspect(c *Config) error {
	if len(c.args) == 0 {
		return errors.New("Not enough arguements\nUsage: inspect <pokemon-name>")
	}

	if len(c.args) > 1 {
		return errors.New("Too many arguements\nUsage: inspect <pokemon-name>")
	}

	pokemon, ok := c.pokedex[c.args[0]]

	if !ok {
		return errors.New("You have not caught that pokemon\n")
	}

	pokemon.Print()
	return nil
}

func commandPokedex(c *Config) error {
    if len(c.pokedex) == 0 {
        return errors.New("You have caught no pokemon. Try catching them using catch\n")
    }
    fmt.Printf("Your Pokedex:\n")
    for k := range c.pokedex {
        fmt.Printf(" - %v\n", k)
    }
    return nil
}
