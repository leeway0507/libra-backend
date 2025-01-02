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
)

func TestScrap(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	pool := db.ConnectPGPool(cfg.DATABASE_URL, ctx)

	t.Run("book detail", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/scrap/%s/%s", "111005", "9791169850483"), nil)
		req.SetPathValue("libCode", "111005")
		req.SetPathValue("isbn", "9791196411725")
		resp := httptest.NewRecorder()
		HandleScrap(resp, req, pool)

		if resp.Result().StatusCode != 200 {
			b, _ := io.ReadAll(resp.Body)
			log.Printf("resp: %#+v\n", string(b))
			t.Fatal("failed to respond")
		}

	})
}
