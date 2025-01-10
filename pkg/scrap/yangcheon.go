package scrap

import (
	"fmt"
	"io"
	"libra-backend/model"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type yangcheon struct {
	model.Lib
}

func NewYangcheon(isbn, district, libname string) model.LibScrap {
	return &yangcheon{
		Lib: model.Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
	}
}

func (e *yangcheon) Request() (io.ReadCloser, error) {
	url, err := url.Parse("https://lib.yangcheon.or.kr/main/site/search/bookSearch.do")
	if err != nil {
		log.Println(err)
	}
	queryParam := url.Query()
	queryParam.Set("detail", "ok")
	queryParam.Set("cmd_name", "booksearch")
	queryParam.Set("search_type", "detail")
	queryParam.Set("search_isbn_issn", e.Isbn)
	url.RawQuery = queryParam.Encode()

	r, err := http.Get(url.String())
	if err != nil {
		log.Println(err)
	}
	if r.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", r.StatusCode)
		return nil, fmt.Errorf("error status 500")
	}
	return r.Body, nil
}

func (e *yangcheon) ExtractData(body io.ReadCloser) (*[]model.LibBookStatus, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var Books []model.LibBookStatus
	doc.Find("div.book_area").Each(func(i int, s *goquery.Selection) {
		libName := strings.Replace(s.Find(" .tit span").Text(), "[", "", 1)
		libName = strings.Replace(libName, "]", "", 1)
		classNum := strings.ReplaceAll(s.Find(".list_area > dl:nth-child(5) > dd").Text(), "\t", "")
		bookStatusRaw := strings.ReplaceAll(s.Find(".book_status").Text(), "\t", "")
		bookStatus := strings.ReplaceAll(bookStatusRaw, "\n", "")

		Books = append(Books, model.LibBookStatus{
			Isbn:       e.Isbn,
			District:   e.District,
			LibName:    libName,
			ClassNum:   strings.Trim(classNum, " \n"),
			BookStatus: bookStatus,
		})
	})

	if Books == nil {
		return nil, fmt.Errorf("ExtractData : no match data")
	}
	return &Books, nil
}

func (e *yangcheon) GetDistrict() string {
	return e.District
}
func (e *yangcheon) GetIsbn() string {
	return e.Isbn
}
func (e *yangcheon) GetLibName() string {
	return e.LibName
}
