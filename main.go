package main

import (
	"fmt"
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
	fmt.Println("Hello, World!")
}
