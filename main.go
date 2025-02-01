package main

import (
	"bufio"
	"fmt"
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
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			slog.Error("Error scanning input")
			break
		}
		input := scanner.Text()
		clearedInput := cleanInput(input)
		for _, word := range clearedInput {
			fmt.Print("Your command was: ", word)
			break
		}
		fmt.Println()
	}
}
