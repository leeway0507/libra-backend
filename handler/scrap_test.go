package handler

import (
	"context"
	"fmt"
	"io"
	"libra-backend/config"
	"libra-backend/db"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgraph-io/ristretto/v2"
)

func TestScrap(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	pool := db.ConnectPGPool(cfg.DATABASE_URL, ctx)

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e3,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 28, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	t.Run("book detail", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/scrap/%s/%s", "111005", "9791169850483"), nil)
		req.SetPathValue("libCode", "111005")
		req.SetPathValue("isbn", "9791196411725")
		resp := httptest.NewRecorder()
		HandleScrap(resp, req, pool, cache)

		if resp.Result().StatusCode != 200 {
			b, _ := io.ReadAll(resp.Body)
			log.Printf("resp: %#+v\n", string(b))
			t.Fatal("failed to respond")
		}

		response, found := cache.Get(req.URL.RawPath)
		if !found {
			t.Fatal("failed to save cache")
		}

		log.Printf("string(response): %#+v\n", string(response))

	})
}
