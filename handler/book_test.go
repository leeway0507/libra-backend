package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookRequest(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)
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
		HandleBookRequests(resp, req, query)

		if resp.Result().StatusCode != 200 {
			t.Fatal("failed to respond")
		}

	})
}
