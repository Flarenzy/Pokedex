package run

import (
	"errors"
	"fmt"
	"os"

	"github.com/Flarenzy/Pokedex/cmd"
	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/chzyer/readline"
)

func Run(c *config.Config, in LineReader, commands map[string]*cmd.CliCommand) error {
	for {
		line, err := in.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				c.Logger.Info("Interrupt received")
				c.Cache.Done()
				os.Exit(0)
			}
			c.Logger.Error(fmt.Sprintf("Error reading line: %s", err))
		}
		clearedInput := cleanInput(line)
		if len(clearedInput) < 1 {
			continue
		}

		command, ok := commands[clearedInput[0]]
		if len(clearedInput) > 1 {
			var args []string
			for _, arg := range clearedInput[1:] {
				args = append(args, arg)
			}
			c.Args = args
		}
		if !ok {
			continue
		}
		err = command.Callback(c)
		if err != nil {
			if errors.Is(err, cmd.ErrStop) {
				c.Cache.Done()
				c.Logger.Info("Exit received")
				return nil
			}
			c.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
