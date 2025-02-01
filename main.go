package main

import (
	"bufio"
	"fmt"
	"github.com/Flarenzy/Pokedex/cmd"
	"log/slog"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	res := make([]string, 0)
	temp := strings.Fields(text)
	for _, v := range temp {
		res = append(res, strings.ToLower(strings.TrimSpace(v)))
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := cmd.NewCommands()
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			slog.Error("Error scanning input")
			break
		}
		input := scanner.Text()
		clearedInput := cleanInput(input)
		if len(clearedInput) > 0 {
			command, ok := commands[clearedInput[0]]
			if !ok {
				continue
			}
			err := command.Callback()
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		}

		fmt.Println()
	}
}
