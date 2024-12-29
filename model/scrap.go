package model

import "io"

type LibScrap interface {
	Request() (io.ReadCloser, error)
	ExtractData(body io.ReadCloser) *[]LibBookStatus
	GetDistrict() string
	GetIsbn() string
}
type LocalScrap interface {
	LibScrap
	SaveReqToLocal()
	ExtractDataFromLocal() *[]LibBookStatus
}
type Lib struct {
	Isbn     string
	District string
	LibName  string
}

type LibBookStatus struct {
	District   string `json:"district"`
	LibName    string `json:"libName"`
	Isbn       string `json:"isbn"`
	BookCode   string `json:"bookCode"`
	BookStatus string `json:"bookStatus"`
}
