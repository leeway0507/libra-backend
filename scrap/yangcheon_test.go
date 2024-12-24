package scrap

import (
	"log"
	"testing"
)

func TestYancheon(t *testing.T) {
	isbn := "8970126740"
	y := NewYangcheon(isbn)
	l := NewLocalTest(y)

	t.Run("request isbn", func(t *testing.T) {
		l.SaveReqToLocal()
	})

	t.Run("load body", func(t *testing.T) {
		r := l.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})

}
