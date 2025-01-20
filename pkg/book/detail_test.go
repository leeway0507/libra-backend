package book

import (
	"context"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"testing"
)

func TestBookDetail(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)
	t.Run("get data by isbn", func(t *testing.T) {
		isbn := "9791169850483"
		libCode := []string{"111005", "111015"}
		_, err := GetBookDetail(query, isbn, libCode)
		if err != nil {
			t.Fatal(err)
		}

	})
}
