package scrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/model"
	"log"
	"net/http"
	"net/url"
)

type RequestBody struct {
	PubFormCode     string   `json:"pubFormCode"`
	Display         string   `json:"display"`
	Article         string   `json:"article"`
	Order           string   `json:"order"`
	ManageCode      []string `json:"manageCode"`
	AdvIsbn         string   `json:"advIsbn"`
	AdvContentsType []string `json:"advContentsType"`
	AdvTextLang     string   `json:"advTextLang"`
}

type Book struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	CoverUrl       string `json:"coverUrl"`
	Isbn           string `json:"isbn"`
	WorkingStatus  string `json:"workingStatus"`
	ReturnPlanDate string `json:"returnPlanDate"`
	CallNo         string `json:"callNo"`
	LibName        string `json:"libName"`
}

type BookResponse struct {
	Contents struct {
		BookList []Book `json:"bookList"`
	} `json:"contents"`
}

type kangseo struct {
	model.Lib
}

func NewKangSeo(isbn, district, libname string) model.LibScrap {
	return &kangseo{
		Lib: model.Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
	}
}

func (e *kangseo) Request() (io.ReadCloser, error) {
	urlRaw, err := url.JoinPath("https://lib.gangseo.seoul.kr/api/search")
	if err != nil {
		log.Println(err)
	}

	body := RequestBody{
		PubFormCode:     "ALL",
		Display:         "10",
		Article:         "SCORE",
		Order:           "DESC",
		ManageCode:      []string{"AG", "BG", "AA", "AB", "AC", "AD", "AF", "AE", "BK", "AQ", "AL", "AR", "AJ", "AI", "AX", "AO", "BJ", "AK", "BB", "BL", "AS", "AN", "BC", "BI", "AZ", "AW", "BD", "AP", "BH", "BF", "AM", "AY", "AU", "AV", "BA", "AT"},
		AdvIsbn:         e.Isbn,
		AdvContentsType: []string{"ALL"},
		AdvTextLang:     "ALL",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("fail to marshal %v", err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", urlRaw, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	// 헤더 추가
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "ko,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", "isNotTodayPopup112=1")
	req.Header.Set("Host", "lib.gangseo.seoul.kr")
	req.Header.Set("Origin", "https://lib.gangseo.seoul.kr")
	req.Header.Set("Referer", "https://lib.gangseo.seoul.kr/DetailSearchResult")
	req.Header.Set("Sec-CH-UA", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// HTTP 클라이언트 생성 및 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("status", resp.StatusCode)
		log.Println("Error sending request:", err)
		return nil, err
	}

	return resp.Body, nil
}

func (e *kangseo) ExtractData(body io.ReadCloser) (*[]model.LibBookStatus, error) {

	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var bookList BookResponse
	err = json.Unmarshal(b, &bookList)
	if err != nil {
		return nil, err
	}
	var Books []model.LibBookStatus
	for _, book := range bookList.Contents.BookList {
		var bookStatus string = book.WorkingStatus
		if bookStatus == "대출중" {
			bookStatus = bookStatus + "(반납예정일:" + book.ReturnPlanDate + ")"
		}
		Books = append(Books, model.LibBookStatus{
			Isbn:       e.Isbn,
			District:   e.District,
			LibName:    book.LibName,
			ClassNum:   book.CallNo,
			BookStatus: bookStatus,
		})
	}

	if Books == nil {
		return nil, fmt.Errorf("ExtractData : no match data")
	}
	return &Books, nil
}

func (e *kangseo) GetDistrict() string {
	return e.District
}
func (e *kangseo) GetIsbn() string {
	return e.Isbn
}
func (e *kangseo) GetLibName() string {
	return e.LibName
}
