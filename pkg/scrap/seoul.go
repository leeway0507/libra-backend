package scrap

import (
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/model"
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
	IsUseBookShelf    string      `json:"isUseBookShelf"`
	IsUseMissrepo     string      `json:"isUseMissrepo"`
	ReserveItemCount  int         `json:"reserveItemCount"`
	IsUseNightLoanreq string      `json:"isUseNightLoanreq"`
	IsUseLoanreq      string      `json:"isUseLoanreq"`
	BranchLocations   string      `json:"branchLocations"`
	IsUseDelivery     string      `json:"isUseDelivery"`
	BranchLoanFlag    *string     `json:"branchLoanFlag"` // nullable field
	CtrlNo            string      `json:"ctrlno"`
	NoHolding         []NoHolding `json:"noholding"`
	ReserveFlag       string      `json:"reserveFlag"`
	ReservationCount  int         `json:"reservationCount"`
	CanReserveCount   int         `json:"canReserveCount"`
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
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

type seoul struct {
	model.Lib
}

func NewSeoul(isbn, district, libname string) model.LibScrap {
	return &seoul{
		Lib: model.Lib{
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

	r, err := http.Get(url.String())
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

func (e *seoul) ExtractData(body io.ReadCloser) (*[]model.LibBookStatus, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	var Books []model.LibBookStatus
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

	Books = append(Books, model.LibBookStatus{
		Isbn:       e.Isbn,
		District:   e.District,
		LibName:    book.Find("a").Text(),
		BookCode:   book.First().Next().Next().Next().Text(),
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
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", userAgent)

	req.Header.Add("Referer", fmt.Sprintf("https://lib.seoul.go.kr/main/searchBrief?st=KWRD&si=TOTAL&q=%s&lmtsn=000000000018&lmtst=OR&lmt0=seoulm", isbn))
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: "G6LnfKrN1Hpt8jLQ8et3fo72DfX13h6JmQlhSSAUQNThobLdvzRSE81ReBHhF0Oa.replibwas_servlet_engine6",
	})
	req.AddCookie(&http.Cookie{
		Name:  "WL_PCID",
		Value: "17361308347406656535282",
	})

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
