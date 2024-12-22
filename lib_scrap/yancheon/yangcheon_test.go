package lib_scrap

import (
	"log"
	"testing"
)

func TestYancheon(t *testing.T) {
	isbn := "8970126740"
	e := NewYangcheon(isbn)

	t.Run("request isbn", func(t *testing.T) {
		e.saveReqToLocal()
	})

	t.Run("load body", func(t *testing.T) {
		r := e.ExtractDataFromLocal()
		log.Printf("r: %#+v\n", r)
	})

}
