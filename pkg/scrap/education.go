package scrap

import (
	"bytes"
	"fmt"
	"io"
	"libra-backend/model"
	"libra-backend/utils"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type education struct {
	model.Lib
	targetLib LibInfo
}
type LibInfo struct {
	libCode  string
	libName  string
	homepage string
	pageName string
}

func NewEduction(isbn, district, libname string) model.LibScrap {
	if !slices.Contains(LibNameArr, libname) {
		return nil
	}

	return &education{
		Lib: model.Lib{
			Isbn:     isbn,
			District: district,
		},
		targetLib: EduLibMap[libname],
	}
}

func (e *education) Request() (io.ReadCloser, error) {
	body := e.searchBook()
	url := e.extractSpecURL(body)

	if url == "" {
		return nil, fmt.Errorf("failed to load book url")
	}
	return e.requestToSpec(url)
}

func (e *education) searchBook() io.ReadCloser {
	urlRaw, err := url.JoinPath(e.targetLib.homepage, "intro/search/index.do")
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}

	// Form Data
	formData := url.Values{
		"collection":    {"new_book"},
		"menu_idx":      {"4"},
		"viewPage":      {"1"},
		"excel_type":    {"SEARCH"},
		"editMode":      {"detail"},
		"locType":       {"library"},
		"searchField":   {"ALL"},
		"isbnCollQuery": {e.Isbn},
		"searchPeriod":  {"m1"},
		"locExquery":    {e.targetLib.libCode},
		"locExquery2":   {e.targetLib.libCode},
		"sortField":     {"RANK"},
		"rowCount":      {"10"},
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", urlRaw, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	// 헤더 추가
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", urlRaw)

	// HTTP 클라이언트 생성 및 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Error sending request:", err)
		return nil
	}

	return resp.Body
}

func (e *education) extractSpecURL(body io.ReadCloser) string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Println("failed to load goquery")
		return ""
	}
	result := doc.Find(".cont .imageType li")
	if result.Length() == 0 {
		log.Println("No Matched ISBN")
		return ""
	}

	vctrl, isExist := result.First().Find("a").First().Attr("vctrl")
	if !isExist {
		log.Println("No Matched vctrl")
		return ""
	}
	urlStr, err := url.JoinPath(e.targetLib.homepage, "intro/search/detail.do")
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("err: %#+v\n", err)
	}
	queryParam := url.Query()
	queryParam.Set("vLoca", e.targetLib.libCode)
	queryParam.Set("vCtrl", vctrl)
	queryParam.Set("isbn", e.Isbn)
	url.RawQuery = queryParam.Encode()

	return url.String()
}

func (e *education) requestToSpec(url string) (io.ReadCloser, error) {
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error status 500")
	}
	if r.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", r.StatusCode)
		return nil, err
	}
	return r.Body, nil

}

func (e *education) ExtractData(body io.ReadCloser) (*[]model.LibBookStatus, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
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
		Isbn:       e.Isbn,
		District:   e.District,
		LibName:    e.targetLib.libName,
		BookCode:   strings.Trim(bookCode, " \n"),
		BookStatus: bookStatus,
	})

	otherLibs := doc.Find("table.otherLibrary > tbody > tr")
	otherLibs.Each(func(i int, s *goquery.Selection) {
		LibName := s.Find("td:nth-child(1)").First().Text()
		bookCode := s.Find("td:nth-child(3)").First().Text()
		bookStatus := s.Find("td:nth-child(5)").First().Text()
		bookBorrowDate := s.Find("td:nth-child(6)").First().Text()
		if bookBorrowDate != " " {
			bookStatus = fmt.Sprintf("%s(%s)", bookStatus, bookBorrowDate)
		}

		Books = append(Books, model.LibBookStatus{
			Isbn:       e.Isbn,
			District:   e.District,
			LibName:    "서울특별시교육청" + LibName,
			BookCode:   bookCode,
			BookStatus: bookStatus,
		})
	})

	if Books == nil {
		return nil, fmt.Errorf("ExtractData : no match data")
	}
	return &Books, nil
}

func (e *education) GetDistrict() string {
	return e.District
}
func (e *education) GetIsbn() string {
	return e.Isbn
}
func (e *education) GetLibName() string {
	return e.LibName
}
