package scrap

import (
	"log"
	"testing"
)

func TestKangnam(t *testing.T) {
	isbn, district, libname := "9791163034735", "강남구", "역삼도서관"
	y := NewKangnam(isbn, district, libname)
	l := NewLocalTest(y)

	t.Run("request isbn", func(t *testing.T) {
		l.SaveReqToLocal()
	})

	t.Run("load body", func(t *testing.T) {
		r, err := l.ExtractDataFromLocal()
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("r: %#+v\n", r)
	})

}
