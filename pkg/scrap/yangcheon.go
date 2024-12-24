package scrap

import (
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

func NewYangcheon(isbn string) model.LibScrap {
	return &yangcheon{
		Lib: model.Lib{
			Isbn:    isbn,
			LibType: "yangcheon",
		},
	}
}

func (e *yangcheon) Request() io.ReadCloser {
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

func (e *yangcheon) ExtractData(body io.ReadCloser) *[]model.LibBookStatus {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	var Books []model.LibBookStatus
	doc.Find("div.book_area").Each(func(i int, s *goquery.Selection) {
		libName := strings.Replace(s.Find(" .tit span").Text(), "[", "", 1)
		libName = strings.Replace(libName, "]", "", 1)
		bookCode := strings.ReplaceAll(s.Find(".list_area > dl:nth-child(5) > dd").Text(), "\t", "")
		bookStatusRaw := strings.ReplaceAll(s.Find(".book_status").Text(), "\t", "")
		bookStatus := strings.ReplaceAll(bookStatusRaw, "\n", "")

		Books = append(Books, model.LibBookStatus{
			Isbn:       "1234",
			LibType:    "yangcheon",
			LibName:    libName,
			BookCode:   strings.Trim(bookCode, " \n"),
			BookStatus: bookStatus,
		})
	})

	// log.Printf("Books: %#+v\n", Books)
	return &Books
}

func (e *yangcheon) GetLibType() string {
	return e.LibType
}
func (e *yangcheon) GetIsbn() string {
	return e.Isbn
}
