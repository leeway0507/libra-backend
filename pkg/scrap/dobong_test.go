package scrap

import (
	"log"
	"testing"
)

func TestDobong(t *testing.T) {
	isbn, district, libname := "9791191590272", "도봉구", "도봉기적의도서관"
	y := NewDobong(isbn, district, libname)
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
