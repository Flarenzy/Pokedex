package run

import (
	"strings"

	"github.com/chzyer/readline"
)

func NewReadlineInput(rl *readline.Instance) *ReadlineInput {
	return &ReadlineInput{rl}
}

func cleanInput(text string) []string {
	res := make([]string, 0)
	temp := strings.Fields(text)
	for _, v := range temp {
		res = append(res, strings.ToLower(strings.TrimSpace(v)))
	}
	return res
}
