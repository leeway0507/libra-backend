package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/kiwi"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	pool := db.ConnectPGPool(cfg.DATABASE_URL, ctx)
	kb := kiwi.NewBuilder(cfg.TOKENIZER_PATH, 1, kiwi.KIWI_BUILD_INTEGRATE_ALLOMORPH)
	k := kb.Build()
	t.Run("no search query", func(t *testing.T) {
		keyword := ""
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/search?", "q=", keyword), nil)
		resp := httptest.NewRecorder()

		HandleSearchQuery(resp, req, pool, k)
		if resp.Result().StatusCode != 400 {
			t.Fatal("should return 400")
		}

	})
	t.Run("search", func(t *testing.T) {
		keywords := []string{"go언어", "go 언어", "go 언어를 활용한"}
		// keywords := []string{"go언어", "go 언어", "go 언어를 활용한", "객체지향", "파이썬", "파이썬 프로그래밍"}

		for _, keyword := range keywords {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/search/normal", "?q=", keyword, "&", "libCode=", "111314"), nil)
			resp := httptest.NewRecorder()

			HandleSearchQuery(resp, req, pool, k)
			if resp.Result().StatusCode != 200 {
				t.Fatal(resp.Body)
			}

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			var books []sqlc.SearchFromBooksRow
			err = json.Unmarshal(b, &books)
			if err != nil {
				t.Fatal(err)
			}
			if len(books) == 0 {
				t.Fatal("len book is 0")
			}
			for _, b := range books[:5] {
				log.Printf("\n %s||%v \n", b.Title.String, b.Score)
			}
		}
	})
}
