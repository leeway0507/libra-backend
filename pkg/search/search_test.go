package search

import (
	"context"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"log"
	"testing"
)

func Test(t *testing.T) {
	cfg := config.GetEnvConfig()
	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)
	search := New(query, cfg.OPEN_AI_API_KEY)
	t.Run("load embedding", func(t *testing.T) {
		_, err := search.LoadEmbeddingQuery("파이썬")
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("not to load embedding", func(t *testing.T) {
		embedding, err := search.LoadEmbeddingQuery("파이썬2")
		log.Println("Intended: ", err)
		if err == nil {
			t.Fatal("should fail to load")
		}
		if embedding != nil {
			t.Fatal("should not return embedding")
		}
	})
}
