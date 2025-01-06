package scrap

import (
	"log"
	"testing"
)

func TestKangseo(t *testing.T) {
	// exist isbn 9791161758428
	isbn, district, libname := "9791161758428", "강서구", "가양도서관"
	y := NewKangSeo(isbn, district, libname)

	t.Run("request isbn", func(t *testing.T) {
		body, err := y.Request()
		if err != nil {
			t.Fatal(err)
		}
		libBookStatus, err := y.ExtractData(body)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("libBookStatus: %#+v\n", libBookStatus)
	})

}
