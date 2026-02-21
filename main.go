package main

import (
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Flarenzy/Pokedex/cmd"
	"github.com/Flarenzy/Pokedex/internal/config"
	internalHTTP "github.com/Flarenzy/Pokedex/internal/http"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
	"github.com/Flarenzy/Pokedex/internal/run"
	"github.com/chzyer/readline"
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
	commands := cmd.NewCommands()
	cache := pokecache.NewCache(50 * time.Second)
	myPokedex := pokedex.NewPokedex()
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex>",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     "/tmp/pokedex.tmp",
	})
	if err != nil {

		panic(err)
	}
	l.CaptureExitSignal()
	rl := run.NewReadlineInput(l)
	defer func() {
		err = l.Close()
		if err != nil {
			slog.Error(err.Error())
		}
		err = rl.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	logLevel := logging.MyHandler{
		Level: slog.LevelDebug,
	}
	logger := logging.NewLogger(logLevel)
	httpClient := internalHTTP.NewDefaultHTTPClient()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := config.NewConfig(cache, logger, myPokedex, httpClient, r.Float64)
	if areaURL := os.Getenv("POKEDEX_AREA_URL"); areaURL != "" {
		c.AreaURL = areaURL
		c.Next = areaURL
	}
	if pokemonURL := os.Getenv("POKEDEX_POKEMON_URL"); pokemonURL != "" {
		c.PokemonURL = pokemonURL
	}
	if err = run.Run(c, rl, commands); err != nil {
		c.Logger.Error(err.Error())
		os.Exit(1)
	}
}
