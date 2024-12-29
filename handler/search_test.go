package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)

	t.Run("no search query", func(t *testing.T) {
		keyword := ""

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/search?", "q=", keyword), nil)
		resp := httptest.NewRecorder()

		HandleSearchNormal(resp, req, query)
		if resp.Result().StatusCode != 400 {
			t.Fatal("should return 400")
		}

	})
	t.Run("search", func(t *testing.T) {
		keyword := "한강"

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/search?", "q=", keyword), nil)
		resp := httptest.NewRecorder()

		HandleSearchNormal(resp, req, query)
		if resp.Result().StatusCode != 200 {
			t.Fatal("failed to respond")
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var books []sqlc.Book
		err = json.Unmarshal(b, &books)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(books[0].Author.String, "한강") {
			log.Printf("books[0]: %#+v\n", books[0])
			t.Fatal("wrong answer")
		}
	})
}
