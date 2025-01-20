package scrap

import (
	"log"
	"testing"
)

func TestSeoul(t *testing.T) {
	// exist isbn : 9791197453649
	// non exist isbn : 9791191590272
	isbn, district, libname := "9791197453649", "서울시", "서울도서관"
	y := NewSeoul(isbn, district, libname)
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
	t.Run("request status", func(t *testing.T) {
		code := "CAT000001643479"
		seoulInstance := &seoul{Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
		}
		status, err := seoulInstance.RequestStatus(isbn, code)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(status)
	})
}
