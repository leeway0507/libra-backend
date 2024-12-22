package model

import "io"

type LibScrap interface {
	Request() io.ReadCloser
	ExtractData() *LibBookStatus
}

type LibBookStatus struct {
	LibType    string
	LibName    string
	Isbn       string
	BookCode   string
	BookStatus string
}
