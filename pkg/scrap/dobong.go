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

type dobong struct {
	model.Lib
}

func NewDobong(isbn, district, libname string) model.BookStatusScraper {
	return &dobong{
		Lib: model.Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
	}
}

func (e *dobong) Request() (io.ReadCloser, error) {
	url, err := url.Parse("https://www.unilib.dobong.kr/site/search/search00.do")
	if err != nil {
		log.Println(err)
	}

	queryParam := url.Query()

	queryParam.Set("detail", "ok")
	queryParam.Set("cmd_name", "bookandnonbooksearch")
	queryParam.Set("search_type", "detail")
	queryParam.Set("search_isbn_issn", e.Isbn)
	queryParam.Set("manage_code", "MA,MB,MC,ME,MG,MJ,MF,MH,SA,MD,SB,SL,SM,SN,SO,SP,SJ,SK,SQ,SS,ST,SU,SG,SH,SC")
	url.RawQuery = queryParam.Encode()

	r, err := http.Get(url.String())
	// log.Println("url.String()", url.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if r.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", r.StatusCode)
		return nil, fmt.Errorf("error status 500")
	}
	return r.Body, nil
}

func (e *dobong) ExtractData(body io.ReadCloser) (*[]model.LibBookStatus, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var Books []model.LibBookStatus
	doc.Find("div.book_area").Each(func(i int, s *goquery.Selection) {
		libName := s.Find(" .tit span.lib_name").Text()
		classNum := strings.ReplaceAll(s.Find(".cont dd").Last().Text(), "\t", "")
		bookStatusRaw := strings.ReplaceAll(s.Find(".book_status").Text(), "\t", "")
		bookStatus := strings.ReplaceAll(bookStatusRaw, "\n", "")
		isEbook := len(s.Find(".ebook_icon").Text()) > 0
		if isEbook {
			bookStatus = "전자책"
			classNum = "-"
		}

		bookReturnDate := s.Find(".cont dl:nth-child(5) dt").Text()
		if bookReturnDate == "반납예정일" {
			date := s.Find(".cont dl:nth-child(5) dd").Text()
			bookStatus = fmt.Sprintf("%s(%s)", bookStatus, date)
		}

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

func (e *dobong) GetDistrict() string {
	return e.District
}
func (e *dobong) GetIsbn() string {
	return e.Isbn
}
func (e *dobong) GetLibName() string {
	return e.LibName
}
