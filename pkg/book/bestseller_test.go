package book

import (
	"encoding/json"
	"libra-backend/config"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	dataPath = "/Users/yangwoolee/repo/libra-backend/data"
)

func TestBestSeller(t *testing.T) {
	cfg := config.GetEnvConfig()
	b := NewBestSeller(cfg.ALADIN_API_KEY)
	t.Run("aladin bestseller", func(t *testing.T) {
		b.FetchBestSellers("1", "0", "2025", "1", "1")
	})
	t.Run("aladin response", func(t *testing.T) {
		f, err := os.Open(filepath.Join(dataPath, "aladin_request.json"))
		if err != nil {
			t.Fatal(err)
		}
		j := b.Parse(f)
		log.Printf("j: %#+v\n", j)
	})
	t.Run("get bestseller default", func(t *testing.T) {
		r := b.GetBestSellerDefault()
		x, _ := json.Marshal(r)
		os.WriteFile(filepath.Join(dataPath, "aladin_response.json"), x, 0600)
	})

}
