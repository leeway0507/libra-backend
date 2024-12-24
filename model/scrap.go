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
	District   string
	LibName    string
	Isbn       string
	BookCode   string
	BookStatus string
}
