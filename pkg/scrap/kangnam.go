package scrap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type kangnam struct {
	Lib
}

var (
	CookieKangnam []*http.Cookie
)

func NewKangnam(isbn, district, libname string) BookStatusScraper {
	return &kangnam{
		Lib: Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
	}
}

func (e *kangnam) Request() (io.ReadCloser, error) {
	urlRaw := "https://library.gangnam.go.kr/intro/menu/10004/program/30002/plusSearchResultList.do"
	// Form Data
	formData := url.Values{
		"searchType":             {"DETAIL"},
		"searchCategory":         {"BOOK"},
		"searchKey1":             {"TITLE"},
		"searchKeyword1":         {""},
		"searchOperator1":        {"AND"},
		"searchKey2":             {"AUTHOR"},
		"searchKeyword2":         {""},
		"searchOperator2":        {"AND"},
		"searchKey3":             {"PUBLISHER"},
		"searchKeyword3":         {""},
		"searchOperator3":        {"AND"},
		"searchKey4":             {"KEYWORD"},
		"searchKeyword4":         {""},
		"searchOperator4":        {"AND"},
		"searchKey5":             {"ISBN"},
		"searchKeyword5":         {e.Isbn},
		"searchOperator5":        {"AND"},
		"searchPublishStartYear": {""},
		"searchPublishEndYear":   {""},
		"searchLibrary":          {"ALL"},
		"searchRoom":             {"ALL"},
		"searchSort":             {"SIMILAR"},
		"searchOrder":            {"DESC"},
		"searchRecordCount":      {"10"},
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", urlRaw, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// 헤더 추가
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", urlRaw)
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Host", "library.gangnam.go.kr")

	for _, c := range CookieKangnam {
		req.AddCookie(c)
	}

	// HTTP 클라이언트 생성 및 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	CookieKangnam = resp.Cookies()
	return resp.Body, nil
}

func (e *kangnam) ExtractData(body io.ReadCloser) (*[]LibBookStatus, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var Books []LibBookStatus
	doc.Find("ul.resultList li").Each(func(i int, s *goquery.Selection) {
		libName, _ := strings.CutPrefix(s.Find("dd.site > span").First().Text(), "도서관: ")
		classNum := strings.ReplaceAll(s.Find("dd.data > span").Last().Text(), "\t", "")
		classNum = strings.ReplaceAll(classNum, "\n", "")
		classNum = strings.ReplaceAll(classNum, "위치출력", "")
		bookStatus := s.Find("b").Text()

		Books = append(Books, LibBookStatus{
			Isbn:       e.Isbn,
			District:   e.District,
			LibName:    libName,
			ClassNum:   classNum,
			BookStatus: bookStatus,
		})
	})

	if Books == nil {
		return nil, fmt.Errorf("ExtractData : no match data")
	}
	return &Books, nil
}

func (e *kangnam) GetDistrict() string {
	return e.District
}
func (e *kangnam) GetIsbn() string {
	return e.Isbn
}
func (e *kangnam) GetLibName() string {
	return e.LibName
}
