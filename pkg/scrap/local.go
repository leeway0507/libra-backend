package scrap

import (
	"io"
	"libra-backend/model"
	"libra-backend/utils"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	model.LibScrap
}

func NewLocalTest(LibScrap model.LibScrap) model.LocalScrap {
	return &Local{
		LibScrap: LibScrap,
	}
}

func (L *Local) SaveReqToLocal() {
	body, err := L.Request()
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	cd := utils.GetCurrentFolderName()
	libType := L.GetDistrict()
	f, err := os.Create(filepath.Join(cd, "test_html", libType+".html"))
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		log.Println(err)
	}

	b, err := io.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	f.Write(b)
}

func (L *Local) ExtractDataFromLocal() *[]model.LibBookStatus {
	cd := utils.GetCurrentFolderName()
	libType := L.GetDistrict()
	f, err := os.Open(filepath.Join(cd, "test_html", libType+".html"))
	if err != nil {
		log.Println(err)
		return nil
	}
	return L.ExtractData(f)
}
