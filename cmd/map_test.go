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

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
)

const locationAreaFixture = `{"count":2,"next":"https://example.test/next","previous":"https://example.test/prev","results":[{"name":"area-one","url":"u1"},{"name":"area-two","url":"u2"}]}`

type readErrBody struct{}

func (r readErrBody) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}

func (r readErrBody) Close() error {
	return nil
}

type closeErrBody struct {
	data []byte
	read bool
}

func (c *closeErrBody) Read(p []byte) (int, error) {
	if c.read {
		return 0, io.EOF
	}
	c.read = true
	n := copy(p, c.data)
	return n, nil
}

func (c *closeErrBody) Close() error {
	return errors.New("close failed")
}

type bodyHTTPClient struct {
	body io.ReadCloser
}

func (b bodyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Body: b.body, Header: make(http.Header), Request: req}, nil
}

func TestCommandMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		next         string
		body         string
		clientErr    error
		out          any
		wantErr      error
		wantErrAny   bool
		wantContains []string
		wantNext     string
		wantPrevious string
	}{
		{name: "success", next: "https://example.test/location-area", body: locationAreaFixture, out: &bytes.Buffer{}, wantContains: []string{"area-one", "area-two"}, wantNext: "https://example.test/next", wantPrevious: "https://example.test/prev"},
		{name: "api error", next: "https://example.test/location-area", clientErr: httpClientDoError, out: &bytes.Buffer{}, wantErr: httpClientDoError},
		{name: "invalid json", next: "https://example.test/location-area", body: "{", out: &bytes.Buffer{}, wantErrAny: true},
		{name: "invalid url", next: "://bad", body: locationAreaFixture, out: &bytes.Buffer{}, wantErrAny: true},
		{name: "location write error", next: "https://example.test/location-area", body: locationAreaFixture, out: errWriter{}, wantErr: writeError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := pokecache.NewCache(20 * time.Second)
			t.Cleanup(cache.Done)

			writer, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{
				Next:       tc.next,
				Cache:      cache,
				Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:        writer,
				HTTPClient: &stubHTTPClient{body: tc.body, err: tc.clientErr},
			}

			err := commandMap(&c)
			if tc.wantErrAny {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if buf, ok := writer.(*bytes.Buffer); ok {
				for _, expected := range tc.wantContains {
					if !strings.Contains(buf.String(), expected) {
						t.Fatalf("expected output to contain %q, got %q", expected, buf.String())
					}
				}
			}
			if tc.wantNext != "" && c.Next != tc.wantNext {
				t.Fatalf("expected next %q, got %q", tc.wantNext, c.Next)
			}
			if tc.wantPrevious != "" && c.Previous != tc.wantPrevious {
				t.Fatalf("expected previous %q, got %q", tc.wantPrevious, c.Previous)
			}
		})
	}
}

func TestCommandMapb(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		previous     string
		body         string
		clientErr    error
		out          any
		wantErr      error
		wantContains string
	}{
		{name: "first page message", previous: "", body: locationAreaFixture, out: &bytes.Buffer{}, wantContains: "you're on the first page"},
		{name: "loads previous page", previous: "https://example.test/previous", body: locationAreaFixture, out: &bytes.Buffer{}, wantContains: "area-one"},
		{name: "write error on first page", previous: "", body: locationAreaFixture, out: errWriter{}, wantErr: writeError},
		{name: "previous page fetch error", previous: "https://example.test/previous", clientErr: httpClientDoError, body: locationAreaFixture, out: &bytes.Buffer{}, wantErr: httpClientDoError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := pokecache.NewCache(20 * time.Second)
			t.Cleanup(cache.Done)

			writer, ok := tc.out.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("invalid writer")
			}
			c := config.Config{
				Previous:   tc.previous,
				Cache:      cache,
				Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:        writer,
				HTTPClient: &stubHTTPClient{body: tc.body, err: tc.clientErr},
			}

			err := commandMapb(&c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if buf, ok := writer.(*bytes.Buffer); ok && tc.wantContains != "" && !strings.Contains(buf.String(), tc.wantContains) {
				t.Fatalf("expected output to contain %q, got %q", tc.wantContains, buf.String())
			}
		})
	}
}

func TestGetFromAPIReadError(t *testing.T) {
	t.Parallel()

	c := config.Config{
		Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
		Out:        &bytes.Buffer{},
		HTTPClient: bodyHTTPClient{body: readErrBody{}},
	}

	_, err := getFromAPI("https://example.test/location-area", &c)
	if err == nil {
		t.Fatal("expected read error")
	}
}

func TestGetFromAPIPanicsOnCloseError(t *testing.T) {
	t.Parallel()

	c := config.Config{
		Logger: logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
		Out:    &bytes.Buffer{},
		HTTPClient: bodyHTTPClient{body: &closeErrBody{
			data: []byte(`{"ok":true}`),
		}},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on close error")
		}
		if r != "error closing body" {
			t.Fatalf("unexpected panic value: %v", r)
		}
	}()

	_, _ = getFromAPI("https://example.test/location-area", &c)
}

func TestNewMapCommands(t *testing.T) {
	t.Parallel()

	mapCmd := newMapCommand()
	if mapCmd == nil || mapCmd.Callback == nil || mapCmd.name != "map" {
		t.Fatal("invalid map command")
	}
	mapbCmd := newMapbCommand()
	if mapbCmd == nil || mapbCmd.Callback == nil || mapbCmd.name != "mapb" {
		t.Fatal("invalid mapb command")
	}
}
