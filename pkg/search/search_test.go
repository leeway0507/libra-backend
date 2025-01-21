package search

import (
	"context"
	"fmt"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
	"log"
	"math"
	"testing"
)

func NumberSliceCosineSimilarity(sliceA, sliceB []float32) float64 {
	r, _ := NumberSliceCosineSimilarityE(sliceA, sliceB)
	return r
}

func NumberSliceCosineSimilarityE(sliceA, sliceB []float32) (float64, error) {
	if len(sliceA) != len(sliceB) {
		return 0, fmt.Errorf("length must equals")
	}
	var t1 float32
	var t2 float32
	var t3 float32
	for index := range sliceA {
		t1 += float32(sliceA[index]) * float32(sliceB[index])
		t2 += float32(sliceA[index]) * float32(sliceA[index])
		t3 += float32(sliceB[index]) * float32(sliceB[index])
	}
	return float64(t1) / (math.Sqrt(float64(t2)) * math.Sqrt(float64(t3))), nil
}

func TestSearch(t *testing.T) {
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
	t.Run("check similarity", func(t *testing.T) {
		e1, _ := search.LoadEmbeddingQuery("파이썬")
		e2, _ := search.LoadEmbeddingQuery("점프 투 파이썬")

		log.Println(NumberSliceCosineSimilarity(e1.Embedding, e2.Embedding))

	})
}
