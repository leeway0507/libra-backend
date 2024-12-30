package handler

import (
	"context"
	"fmt"
	"io"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrap(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)
	t.Run("book detail", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/scrap/%s/%s", "111005", "9791169850483"), nil)
		req.SetPathValue("libCode", "111005")
		req.SetPathValue("isbn", "9791196411725")
		resp := httptest.NewRecorder()
		HandleScrap(resp, req, query)

		if resp.Result().StatusCode != 200 {
			b, _ := io.ReadAll(resp.Body)
			log.Printf("resp: %#+v\n", string(b))
			t.Fatal("failed to respond")
		}

	})
}
