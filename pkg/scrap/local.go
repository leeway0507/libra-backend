package scrap

import (
	"io"
	"libra-backend/pkg/utils"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	BookStatusScraper
}

func NewLocalTest(LibScrap BookStatusScraper) LocalScrap {
	return &Local{
		BookStatusScraper: LibScrap,
	}
}

func (L *Local) SaveReqToLocal() {
	body, err := L.Request()
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	cd := utils.GetCurrentFolderDir()
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

func (L *Local) ExtractDataFromLocal() (*[]LibBookStatus, error) {
	cd := utils.GetCurrentFolderDir()
	libType := L.GetDistrict()
	f, err := os.Open(filepath.Join(cd, "test_html", libType+".html"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	d, err := L.ExtractData(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return d, nil
}
