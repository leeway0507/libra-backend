package scrap

import (
	"io"
	"libra-backend/model"
	"libra-backend/utils"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type education struct {
	model.Lib
}

func NewEduction(isbn string) model.LibScrap {
	return &education{
		Lib: model.Lib{
			Isbn:    isbn,
			LibType: "education",
		},
	}
}

func (e *education) Request() io.ReadCloser {
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
	// log.Println("url.String()", url.String())
	if err != nil {
		log.Println(err)
	}
	if r.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", r.StatusCode)
	}
	return r.Body
}

func (e *education) ExtractData(body io.ReadCloser) *[]model.LibBookStatus {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	bookCode := doc.Find("#position > tbody > tr > td:nth-child(3)").Text()
	bookStatusRaw := doc.Find("#position > tbody > tr > td:nth-child(5)").Text()
	var bookStatus string = bookStatusRaw
	if bookStatusRaw != "대출가능" {
		ss := strings.Split(strings.ReplaceAll(bookStatusRaw, " ", ""), "\n")
		ss2 := utils.RemoveEmptyStringInSlice(ss)
		bookStatus = strings.Join(ss2, " ")
	}
	var Books []model.LibBookStatus
	Books = append(Books, model.LibBookStatus{
		Isbn:       "1234",
		LibType:    "yangcheon",
		LibName:    "서울특별시교육청양천도서관",
		BookCode:   strings.Trim(bookCode, " \n"),
		BookStatus: bookStatus,
	})
	return &Books

}

func (e *education) GetLibType() string {
	return e.LibType
}
func (e *education) GetIsbn() string {
	return e.Isbn
}
