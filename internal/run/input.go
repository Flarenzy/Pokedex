package run

import "github.com/chzyer/readline"

type LineReader interface {
	Readline() (string, error)
	Close() error
}

type ReadlineInput struct {
	rl *readline.Instance
}

func (r *ReadlineInput) Readline() (string, error) {
	return r.rl.Readline()
}

func (r *ReadlineInput) Close() error {
	return r.rl.Close()
}
