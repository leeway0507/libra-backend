package scrap

import (
	"log"
	"testing"
)

func TestYanGcheon(t *testing.T) {
	isbn, district, libname := "8970126740", "양천구", "갈산도서관"
	y := NewYangcheon(isbn, district, libname)
	l := NewLocalTest(y)

	t.Run("request isbn", func(t *testing.T) {
		l.SaveReqToLocal()
	})

	t.Run("load body", func(t *testing.T) {
		r := l.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})

}
