package main

import (
	"bufio"
	"fmt"
	"github.com/Flarenzy/Pokedex/cmd"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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
	cache := pokecache.NewCache(50 * time.Second)
	sigChan := make(chan os.Signal, 1)
	logLevel := logging.MyHandler{
		Level: slog.LevelDebug,
	}
	logger := logging.NewLogger(logLevel)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		// cleanup
		cache.Done()
		os.Exit(0)
	}()
	config := cmd.NewConfig(cache, logger)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			logger.Error("Error scanning input")
			break
		}
		input := scanner.Text()
		clearedInput := cleanInput(input)
		if len(clearedInput) > 0 {
			command, ok := commands[clearedInput[0]]
			if !ok {
				continue
			}
			err := command.Callback(config)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}
		fmt.Println()
	}
}
