package search

import (
	"log"
	"path/filepath"
	"testing"
)

const TEST_PATH = "/Users/yangwoolee/repo/libra-backend/data/search/grammer"

func TestSpacingCheckAPI(t *testing.T) {
	keywords := []string{"go언어", "go 언어", "go언어를 활용한", "go 어너", "go 어너를 활용한 go언어"}
	correctedWords := []string{"go 언어", "go 언어", "go 언어를 활용한", "go 어너", "go 어너를 활용한 go 언어"}

	t.Run("fetch Spacing", func(t *testing.T) {

		for _, keyword := range keywords {
			fetchSpacingCheckToHtml(TEST_PATH, keyword)

			doc := readHtml(filepath.Join(TEST_PATH, keyword+".html"))

			if doc == nil {
				t.Fatal("fail to load doc for the keyword:", keyword)
			}

		}

	})
	t.Run("extract data from response html", func(t *testing.T) {
		keywords := []string{"go 언어", "go 어너"}
		for _, keyword := range keywords {
			doc := readHtml(filepath.Join(TEST_PATH, keyword+".html"))

			data, err := extractResult(doc)
			if err != nil {
				t.Fatal(err)
			}

			if len(data) == 0 {
				t.Fatal("no data from extractResult")
			}
			log.Printf("data: %#+v\n", data)
		}
	})
	t.Run("exctract corrected keyword", func(t *testing.T) {
		for i, keyword := range keywords {
			doc := readHtml(filepath.Join(TEST_PATH, keyword+".html"))

			correctedWord := extractCorrectedWord(keyword, doc)
			if correctedWords[i] != correctedWord {
				t.Fatalf("result is not matched with expect \n result:%v\n expect:%v\n", correctedWords[i], correctedWord)
			}

		}
	})
}
