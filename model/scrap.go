package model

import "io"

type LibScrap interface {
	Request() (io.ReadCloser, error)
	ExtractData(body io.ReadCloser) (*[]LibBookStatus, error)
	GetDistrict() string
	GetIsbn() string
	GetLibName() string
}
type LocalScrap interface {
	LibScrap
	SaveReqToLocal()
	ExtractDataFromLocal() (*[]LibBookStatus, error)
}
type Lib struct {
	Isbn     string
	District string
	LibName  string
}

type LibBookStatus struct {
	District       string `json:"district"`
	LibName        string `json:"libName"`
	Isbn           string `json:"isbn"`
	BookCode       string `json:"bookCode"`
	BookStatus     string `json:"bookStatus"`
	Toc            string `json:"toc"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}
