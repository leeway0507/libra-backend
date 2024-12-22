package lib_scrap

import (
	"log"
	"testing"
)

func TestEducation(t *testing.T) {

	t.Run("load body", func(t *testing.T) {
		e := NewEduction()
		r := e.ExtractData()
		log.Printf("r: %#+v\n", r)
	})

}
