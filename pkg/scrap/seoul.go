package scrap

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StateResponse struct {
	Location Location `json:"location"`
}

type Location struct {
	NoHolding []NoHolding `json:"noholding"`
}

type NoHolding struct {
	PlaceNo                 string `json:"place_no"`
	BookStateCode           string `json:"book_state_code"`
	PrintAccessionNo        string `json:"print_accessionno"`
	Location                string `json:"location"`
	ReqSubLocaYN            string `json:"req_sub_loca_yn"`
	PlaceName               string `json:"place_name"`
	Bookshelf               string `json:"bookshelf"`
	ReturnPlanDate          string `json:"return_plan_date"`
	MsgPeopleReserved       string `json:"msg_people_reserved"`
	MsgPossibleReserve      string `json:"msg_possible_reserve"`
	SubLocation             string `json:"sub_location"`
	CallNo                  string `json:"call_no"`
	PreserveFlag            string `json:"preserveFlag"`
	LoanFlag                string `json:"loan_flag"`
	BookStatus              string `json:"book_status"`
	Reservation             string `json:"reservation"`
	MainNo                  string `json:"main_no"`
	MsgReserveLimitExceeded string `json:"msg_reserve_limit_exceeded"`
	AccessionNo             string `json:"accession_no"`
	BookState               string `json:"book_state"`
	OrgMarcType             string `json:"orgMarcType"`
	ReqStatus               string `json:"req_status"`
}

var (
	RequestCookie_SEOUL       []*http.Cookie
	RequestStatusCookie_SEOUL []*http.Cookie
)

type seoul struct {
	Lib
}

func NewSeoul(isbn, district, libname string) BookStatusScraper {
	return &seoul{
		Lib: Lib{
			Isbn:     isbn,
			District: district,
			LibName:  libname,
		},
	}
}

func (e *seoul) Request() (io.ReadCloser, error) {
	url, err := url.Parse("https://lib.seoul.go.kr/main/searchBrief")
	if err != nil {
		log.Println(err)
	}

	queryParam := url.Query()

	queryParam.Set("st", "KWRD")
	queryParam.Set("si", "TOTAL")
	queryParam.Set("sts", "Y")
	queryParam.Set("q", e.Isbn)
	url.RawQuery = queryParam.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Host", "lib.seoul.go.kr")
	req.Header.Add("Referer", "https://lib.seoul.go.kr/")

	if len(RequestCookie_SEOUL) != 0 {
		for _, c := range RequestCookie_SEOUL {
			req.AddCookie(c)
		}
	}

	// HTTP 클라이언트 생성 및 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Error sending request:", err)
		return nil, nil
	}
	if resp.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", resp.StatusCode)
		return nil, fmt.Errorf("error status 500")
	}

	RequestCookie_SEOUL = resp.Cookies()

	return resp.Body, nil
}

func (e *seoul) ExtractData(body io.ReadCloser) (*[]LibBookStatus, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var Books []LibBookStatus
	book := doc.Find("ul.book-list > li div.info")

	var bookStatus string
	code := doc.Find(".btnStub").AttrOr("id", "")
	log.Printf("code: %s", code)
	if code == "" {
		bookStatus = ""
	} else {
		code, _ = strings.CutPrefix(code, "btnStub_")
		bookStatus, err = e.RequestStatus(e.Isbn, code)
		if err != nil {
			return nil, err
		}
	}

	classNum := book.Find("p:nth-child(4)").Text()

	Books = append(Books, LibBookStatus{
		Isbn:       e.Isbn,
		District:   e.District,
		LibName:    e.LibName,
		ClassNum:   classNum,
		BookStatus: bookStatus,
	})

	if Books == nil {
		return nil, fmt.Errorf("ExtractData : no match data")
	}
	return &Books, nil
}

func (e *seoul) RequestStatus(isbn string, code string) (string, error) {
	rawUrl, err := url.Parse(fmt.Sprintf("https://lib.seoul.go.kr/search/prevLocJson/%s", code))
	if err != nil {
		log.Println(err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", rawUrl.String(), nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", nil
	}
	if len(RequestStatusCookie_SEOUL) != 0 {
		for _, c := range RequestStatusCookie_SEOUL {
			req.AddCookie(c)
		}
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Referer", fmt.Sprintf("https://lib.seoul.go.kr/main/searchBrief?st=KWRD&si=TOTAL&q=%s&lmtsn=000000000018&lmtst=OR&lmt0=seoulm", isbn))

	// HTTP 클라이언트 생성 및 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Error sending request:", err)
		return "", nil
	}
	if resp.StatusCode != 200 {
		log.Printf("r.StatusCode: %#+v\n", resp.StatusCode)
		return "", fmt.Errorf("error status 500")
	}

	RequestStatusCookie_SEOUL = resp.Cookies()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var BookState StateResponse
	err = json.Unmarshal(b, &BookState)
	if err != nil {
		log.Println(err)
		return "", err
	}
	state := BookState.Location.NoHolding[0].BookState

	if state == "대출중" {
		state = state + "(반납예정일:" +
			BookState.Location.NoHolding[0].ReturnPlanDate +
			")"
	}
	return state, nil
}

func (e *seoul) GetDistrict() string {
	return e.District
}
func (e *seoul) GetIsbn() string {
	return e.Isbn
}
func (e *seoul) GetLibName() string {
	return e.LibName
}
