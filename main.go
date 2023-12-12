package main

import (
	"bufio"
	"fmt"
	"internal/pokecache"
	"os"
	"strings"
	"time"
)

func main() {
	commandMap := getCliCommand()
	scanner := bufio.NewScanner(os.Stdin)
	dur, _ := time.ParseDuration("10s")
	config := Config{
		cache:   pokecache.NewCache(dur),
		pokedex: make(map[string]Pokemon),
	}
	for {
		input := getCommand(scanner)
		inputArr := strings.Split(input, " ")
		if len(inputArr) > 1 {
			config.args = inputArr[1:]
		} else {
			config.args = nil
		}
		command, ok := commandMap[inputArr[0]]
		if !ok {
			fmt.Printf(`Error: %s is not a valid command

            Use help to see list of commands`, input)
			fmt.Printf("\n\n")
			continue
		}
		err := command.callback(&config)
		if err != nil {
			fmt.Printf("Error executing %s:\n%s", input, err)
		}
		fmt.Printf("\n")
	}
}
