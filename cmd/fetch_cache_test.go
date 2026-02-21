package cmd

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/Flarenzy/Pokedex/internal/config"
	"github.com/Flarenzy/Pokedex/internal/domain"
	"github.com/Flarenzy/Pokedex/internal/logging"
	"github.com/Flarenzy/Pokedex/internal/pokecache"
)

func TestGetBodyWithCache(t *testing.T) {
	t.Parallel()

	cacheMissErr := errors.New("cache miss")
	body := []byte(`{"ok":true}`)
	addErr := errors.New("cache add failed")
	url := "https://example.test/resource"

	tests := []struct {
		name       string
		cache      domain.Cacher
		httpClient stubHTTPClient
		wantErr    error
		wantBody   []byte
		wantErrAny bool
	}{
		{
			name: "cache hit",
			cache: &stubCache{
				getBody: body,
			},
			httpClient: stubHTTPClient{err: httpClientDoError},
			wantBody:   body,
		},
		{
			name: "cache miss then api success",
			cache: &stubCache{
				getErr: cacheMissErr,
			},
			httpClient: stubHTTPClient{body: string(body)},
			wantBody:   body,
		},
		{
			name: "cache miss then api error",
			cache: &stubCache{
				getErr: cacheMissErr,
			},
			httpClient: stubHTTPClient{err: httpClientDoError},
			wantErr:    httpClientDoError,
		},
		{
			name: "cache miss then add key exists",
			cache: &stubCache{
				getErr: cacheMissErr,
				addErr: pokecache.ErrKeyExists,
			},
			httpClient: stubHTTPClient{body: string(body)},
			wantBody:   body,
		},
		{
			name: "cache miss then add unexpected error",
			cache: &stubCache{
				getErr: cacheMissErr,
				addErr: addErr,
			},
			httpClient: stubHTTPClient{body: string(body)},
			wantErr:    addErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := config.Config{
				Cache:      tc.cache,
				Logger:     logging.NewLogger(logging.MyHandler{Level: slog.LevelError}),
				Out:        &bytes.Buffer{},
				HTTPClient: tc.httpClient,
			}

			got, err := getBodyWithCache(&c, url)
			if tc.wantErrAny {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if tc.wantErr == nil && !bytes.Equal(got, tc.wantBody) {
				t.Fatalf("expected body %q, got %q", string(tc.wantBody), string(got))
			}
		})
	}
}
