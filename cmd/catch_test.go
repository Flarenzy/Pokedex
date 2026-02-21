package cmd

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Flarenzy/Pokedex/internal"
	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
	"github.com/Flarenzy/Pokedex/internal/pokedex"
)

var (
	writeError         = errors.New("forced write error")
	httpClientDoError  = errors.New("forced http client do error")
	trailingWriteError = errors.New("forced trailing write error")
	cacheAddError      = errors.New("forced cache add error")
)

type errWriter struct{}

func (e errWriter) Write(p []byte) (n int, err error) {
	return 0, writeError
}

type failOnWriteN struct {
	n     int
	count int
	err   error
}

func (w *failOnWriteN) Write(p []byte) (int, error) {
	w.count++
	if w.count == w.n {
		return 0, w.err
	}
	return len(p), nil
}

type stubHTTPClient struct {
	body string
	err  error
}

func (s stubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(s.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type stubCache struct {
	getBody []byte
	getErr  error
	addErr  error
}

func (s *stubCache) Get(key string) ([]byte, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.getBody, nil
}

func (s *stubCache) Add(key string, val []byte) error {
	return s.addErr
}

func (s *stubCache) Done() {}

const mewFixture = `{
  "id": 151,
  "name": "mew",
  "base_experience": 100,
  "height": 4,
  "is_default": true,
  "order": 215,
  "weight": 40,
  "abilities": [],
  "stats": [],
  "types": []
}`

func TestGetPokemon(t *testing.T) {
	t.Parallel()

	url := "https://example.test/api/v2/pokemon/mew"
	tests := []struct {
		name          string
		rand          float64
		outFactory    func() io.Writer
		body          string
		clientErr     error
		wantErr       error
		wantErrAny    bool
		wantOutput    string
		wantInPokedex bool
	}{
		{
			name:          "caught",
			rand:          0.0,
			outFactory:    func() io.Writer { return &bytes.Buffer{} },
			body:          mewFixture,
			wantOutput:    "mew was caught!",
			wantInPokedex: true,
		},
		{
			name:          "escaped",
			rand:          0.99,
			outFactory:    func() io.Writer { return &bytes.Buffer{} },
			body:          mewFixture,
			wantOutput:    "mew escaped!",
			wantInPokedex: false,
		},
		{
			name:          "write error on caught line",
			rand:          0.0,
			outFactory:    func() io.Writer { return errWriter{} },
			body:          mewFixture,
			wantErr:       writeError,
			wantInPokedex: true,
		},
		{
			name:          "write error on escaped line",
			rand:          0.99,
			outFactory:    func() io.Writer { return errWriter{} },
			body:          mewFixture,
			wantErr:       writeError,
			wantInPokedex: false,
		},
		{
			name:          "api error from getFromAPI",
			rand:          0.0,
			outFactory:    func() io.Writer { return &bytes.Buffer{} },
			body:          mewFixture,
			clientErr:     httpClientDoError,
			wantErr:       httpClientDoError,
			wantInPokedex: false,
		},
		{
			name:          "invalid json response",
			rand:          0.0,
			outFactory:    func() io.Writer { return &bytes.Buffer{} },
			body:          "{",
			wantErrAny:    true,
			wantInPokedex: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := pokecache.NewCache(20 * time.Second)
			t.Cleanup(cache.Done)

			p := pokedex.NewPokedex()
			out := tc.outFactory()
			c := config.Config{
				Pokedex: p,
				Cache:   cache,
				Logger: logging.NewLogger(logging.MyHandler{
					Level: slog.LevelError,
				}),
				Out:         out,
				HTTPClient:  &stubHTTPClient{body: tc.body, err: tc.clientErr},
				RandFloat64: func() float64 { return tc.rand },
			}

			err := getPokemon(&c, url)
			if tc.wantErrAny {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}

			buf, ok := out.(*bytes.Buffer)
			if ok && tc.wantOutput != "" && !strings.Contains(buf.String(), tc.wantOutput) {
				t.Fatalf("expected output to contain %q, got: %q", tc.wantOutput, buf.String())
			}

			_, err = p.GetPokemonByName("mew")
			if tc.wantInPokedex && err != nil {
				t.Fatalf("expected pokemon to be added to pokedex, got error: %v", err)
			}
			if !tc.wantInPokedex && err == nil {
				t.Fatal("expected pokemon not to be added to pokedex")
			}
		})
	}
}

func TestGetPokemonCacheAddHandling(t *testing.T) {
	t.Parallel()

	url := "https://example.test/api/v2/pokemon/mew"
	cacheMissErr := errors.New("cache miss")
	tests := []struct {
		name       string
		addErr     error
		wantErr    error
		wantOutput string
	}{
		{
			name:       "duplicate key error is ignored",
			addErr:     pokecache.ErrKeyExists,
			wantErr:    nil,
			wantOutput: "mew was caught!",
		},
		{
			name:    "unexpected add error is returned",
			addErr:  cacheAddError,
			wantErr: cacheAddError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p := pokedex.NewPokedex()
			c := config.Config{
				Pokedex: p,
				Cache: &stubCache{
					getErr: cacheMissErr,
					addErr: tc.addErr,
				},
				Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:    out,
				HTTPClient: &stubHTTPClient{
					body: mewFixture,
				},
				RandFloat64: func() float64 { return 0.0 },
			}

			err := getPokemon(&c, url)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if tc.wantOutput != "" && !strings.Contains(out.String(), tc.wantOutput) {
				t.Fatalf("expected output to contain %q, got: %q", tc.wantOutput, out.String())
			}
		})
	}
}

func TestCommandCatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		args               []string
		outFactory         func() io.Writer
		client             *stubHTTPClient
		rand               float64
		wantErr            error
		wantOutputContains []string
	}{
		{
			name:       "no pokemon args",
			args:       []string{},
			outFactory: func() io.Writer { return &bytes.Buffer{} },
			wantErr:    ErrNoPokemon,
		},
		{
			name:       "throwing pokeball write error",
			args:       []string{"mew"},
			outFactory: func() io.Writer { return errWriter{} },
			wantErr:    writeError,
		},
		{
			name:       "getPokemon error is logged and command continues",
			args:       []string{"mew"},
			outFactory: func() io.Writer { return &bytes.Buffer{} },
			client:     &stubHTTPClient{body: "{"},
			rand:       0.0,
			wantErr:    nil,
			wantOutputContains: []string{
				"Throwing a Pokeball at mew...",
			},
		},
		{
			name:       "trailing newline write error",
			args:       []string{"mew"},
			outFactory: func() io.Writer { return &failOnWriteN{n: 3, err: trailingWriteError} },
			client:     &stubHTTPClient{body: mewFixture},
			rand:       0.0,
			wantErr:    trailingWriteError,
		},
		{
			name:       "successful catch flow",
			args:       []string{"mew"},
			outFactory: func() io.Writer { return &bytes.Buffer{} },
			client:     &stubHTTPClient{body: mewFixture},
			rand:       0.0,
			wantErr:    nil,
			wantOutputContains: []string{
				"Throwing a Pokeball at mew...",
				"mew was caught!",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := pokecache.NewCache(20 * time.Second)
			t.Cleanup(cache.Done)

			out := tc.outFactory()
			c := config.Config{
				PokemonURL: internal.SecondURL,
				Args:       tc.args,
				Pokedex:    pokedex.NewPokedex(),
				Cache:      cache,
				Logger: logging.NewLogger(logging.MyHandler{
					Level: slog.LevelError,
				}),
				Out:         out,
				HTTPClient:  tc.client,
				RandFloat64: func() float64 { return tc.rand },
			}

			err := commandCatch(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}

			buf, ok := out.(*bytes.Buffer)
			if ok {
				for _, expected := range tc.wantOutputContains {
					if !strings.Contains(buf.String(), expected) {
						t.Fatalf("expected output to contain %q, got: %q", expected, buf.String())
					}
				}
			}
		})
	}
}

func TestGetPokemonUsesCacheWhenPresent(t *testing.T) {
	t.Parallel()

	cache := pokecache.NewCache(20 * time.Second)
	t.Cleanup(cache.Done)

	out := &bytes.Buffer{}
	p := pokedex.NewPokedex()
	url := "https://example.test/api/v2/pokemon/mew"
	err := cache.Add(url, []byte(mewFixture))
	if err != nil {
		t.Fatalf("expected cache.Add to succeed, got %v", err)
	}

	c := config.Config{
		Pokedex: p,
		Cache:   cache,
		Logger: logging.NewLogger(logging.MyHandler{
			Level: slog.LevelError,
		}),
		Out:         out,
		HTTPClient:  &stubHTTPClient{err: errors.New("http should not be called")},
		RandFloat64: func() float64 { return 0.0 },
	}

	err = getPokemon(&c, url)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "mew was caught!") {
		t.Fatalf("expected output to contain catch message, got: %q", out.String())
	}
}

func TestNewCatchCommand(t *testing.T) {
	t.Parallel()

	command := newCatchCommand()
	if command == nil {
		t.Fatal("expected command not to be nil")
	}
	if command.name != "catch" {
		t.Fatalf("expected command name to be catch, got %q", command.name)
	}
	if command.description == "" {
		t.Fatal("expected command description not to be empty")
	}
	if command.Callback == nil {
		t.Fatal("expected command callback not to be nil")
	}
}
