package scrap

import (
	"log"
	"testing"
)

func TestEducation(t *testing.T) {

	isbn, district, libname := "8970126740", "교육청", "서울특별시교육청강서도서관"
	e := NewEduction(isbn, district, libname)
	t.Run("check target Lib", func(t *testing.T) {
		if e == nil {
			t.Fatal("failed to load new education Instance")
		}
	})
	l := NewLocalTest(e)

	t.Run("request isbn", func(t *testing.T) {
		l.SaveReqToLocal()
	})

	t.Run("load body", func(t *testing.T) {
		r := l.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})

}
