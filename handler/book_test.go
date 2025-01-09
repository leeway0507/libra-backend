package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/pkg/book"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/dgraph-io/ristretto/v2"
)

func TestBookRequest(t *testing.T) {
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
		mockBody := &DetailRequest{
			Isbn:     "9791169850483",
			LibCodes: []string{"111005", "111015"},
		}
		b, err := json.Marshal(mockBody)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, "/books/detail", bytes.NewReader(b))
		resp := httptest.NewRecorder()
		HandleBookRequests(resp, req, pool)

		if resp.Result().StatusCode != 200 {
			t.Fatal("failed to respond")
		}
	})
	t.Run("best seller", func(t *testing.T) {
		var bestSellerMock book.BestSellerFn = func() book.AladinResponse {
			b, err := os.ReadFile("/Users/yangwoolee/repo/libra-backend/data/aladin_response.json")
			if err != nil {
				t.Fatal(err)
			}
			var items []book.AladinItem
			err = json.Unmarshal(b, &items)
			if err != nil {
				t.Fatal(err)
			}
			return book.AladinResponse{
				CatValue: "0",
				CatName:  "all",
				Items:    items,
			}
		}

		req, _ := http.NewRequest(http.MethodGet, "/books/bestseller/all?libCode=111005,111015", nil)
		resp := httptest.NewRecorder()
		HandleBestSellerRequests(resp, req, "all", bestSellerMock, pool, cache)
		if resp.Result().StatusCode != 200 {
			t.Fatal("failed to respond")
		}

		var j book.AladinResponse
		b, err := io.ReadAll(resp.Result().Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(b, &j)
		if err != nil {
			t.Fatal(err)
		}

		if len(j.Items) == 0 {
			t.Fatal("response should not be 0")
		}

		result, isCached := cache.Get("all 111005,111015")
		if isCached {
			t.Fatal("cache should exist")
		}

		if reflect.DeepEqual(j, result) {
			t.Fatal("cache should be hit")
		}
	})
}
