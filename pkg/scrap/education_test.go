package scrap

import (
	"log"
	"testing"
)

func TestEducation(t *testing.T) {
	isbn := "8970126740"
	e := NewEduction(isbn)
	l := NewLocalTest(e)

	t.Run("save to local", func(t *testing.T) {
		r := l.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})
	t.Run("load body", func(t *testing.T) {
		r := l.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})

}
