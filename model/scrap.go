package model

import "io"

type LibScrap interface {
	Request() io.ReadCloser
	ExtractData(body io.ReadCloser) *[]LibBookStatus
	GetLibType() string
	GetIsbn() string
}
type LocalScrap interface {
	LibScrap
	SaveReqToLocal()
	ExtractDataFromLocal() *[]LibBookStatus
}
type Lib struct {
	Isbn    string
	LibType string
}
type LibBookStatus struct {
	LibType    string
	LibName    string
	Isbn       string
	BookCode   string
	BookStatus string
}
