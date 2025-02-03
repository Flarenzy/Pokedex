package main

import (
	"bufio"
	//"bytes"
	"fmt"
	"github.com/Flarenzy/Pokedex/cmd"
	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
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
	myPokedex := pokedex.NewPokedex()

	sigChan := make(chan os.Signal, 1)
	logLevel := logging.MyHandler{
		Level: slog.LevelDebug,
	}
	logger := logging.NewLogger(logLevel)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	//usedCommands := []string{}
	//curCommand := 0
	//wg := &sync.WaitGroup{}
	go func() {
		<-sigChan
		// cleanup
		cache.Done()
		os.Exit(0)
	}()
	c := config.NewConfig(cache, logger, myPokedex)
	for {
		// ^[[A - up
		// ^[[B - down
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			logger.Error("Error scanning input")
			break
		}

		input := scanner.Text()
		// TODO: basic console up and down
		//if bytes.Equal([]byte(input), []byte("^[[A")) {
		//	if curCommand-1 > 0 {
		//		fmt.Print("\r \r")
		//		fmt.Print(usedCommands[curCommand-1])
		//		curCommand -= 1
		//	}
		//	continue
		//} else if bytes.Equal([]byte(input), []byte("^[[B")) {
		//	if curCommand+1 < len(usedCommands) {
		//		fmt.Print("\r \r")
		//		fmt.Print(usedCommands[curCommand+1])
		//		curCommand += 1
		//	}
		//	continue
		//}

		clearedInput := cleanInput(input)
		if len(clearedInput) > 0 {
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
			err := command.Callback(c)
			if err != nil {
				logger.Error(err.Error())
			}
		}
		fmt.Println()
		//usedCommands = append(usedCommands, input)
		//curCommand = len(usedCommands) - 1
	}
}
