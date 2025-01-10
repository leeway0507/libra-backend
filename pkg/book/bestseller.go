package book

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type AladinItems struct {
	Item []AladinItem `json:"item"`
}

type AladinItem struct {
	Title        string   `json:"title"`
	Link         string   `json:"link"`
	Author       string   `json:"author"`
	PubDate      string   `json:"pubDate"`
	Description  string   `json:"description"`
	Isbn13       string   `json:"isbn13"`
	ItemId       int      `json:"itemId"`
	Cover        string   `json:"cover"`
	CategoryId   int      `json:"categoryId"`
	CategoryName string   `json:"categoryName"`
	Publisher    string   `json:"publisher"`
	BestDuration string   `json:"bestDuration"`
	BestRank     int      `json:"bestRank"`
	LibCode      []string `json:"libCode"`
}

type AladinResponse struct {
	CatValue string       `json:"catValue"`
	CatName  string       `json:"catName"`
	Items    []AladinItem `json:"items"`
}

type BookCat struct {
	value   string
	korName string
	engName string
}

var (
	aladinItemUrl = "http://www.aladin.co.kr/ttb/api/ItemList.aspx"
	cat           = &[]BookCat{
		{value: "0", korName: "분야 선택", engName: "all"},
		{value: "55890", korName: "건강/취미", engName: "health"},
		{value: "656", korName: "인문학", engName: "humanity"},
		{value: "170", korName: "경제경영", engName: "business"},
		{value: "1", korName: "소설/시/희곡", engName: "novel"},
		{value: "1196", korName: "여행", engName: "travel"},
		{value: "74", korName: "역사", engName: "history"},
		{value: "336", korName: "자기계발", engName: "selfdev"},
		{value: "351", korName: "컴퓨터/모바일", engName: "dev"},
		// {value: "2105", korName: "고전"},
		// {value: "987", korName: "과학"},
		// {value: "4395", korName: "달력/기타"},
		// {value: "8257", korName: "대학교재/전문서적"},
		// {value: "2551", korName: "만화"},
		// {value: "798", korName: "사회과학"},
		// {value: "1383", korName: "수험서/자격증"},
		// {value: "1108", korName: "어린이"},
		// {value: "55889", korName: "에세이"},
		// {value: "517", korName: "예술/대중문화"},
		// {value: "1322", korName: "외국어"},
		// {value: "1230", korName: "요리/살림"},
		// {value: "13789", korName: "유아"},
		// {value: "2913", korName: "잡지"},
		// {value: "112011", korName: "장르소설"},
		// {value: "17195", korName: "전집/중고전집"},
		// {value: "1237", korName: "종교/역학"},
		// {value: "2030", korName: "좋은부모"},
		// {value: "1137", korName: "청소년"},
	}
)

type BestSellerFn = func() AladinResponse

type bestSeller struct {
	ttbKey string
}

func NewBestSeller(ttbKey string) *bestSeller {
	return &bestSeller{
		ttbKey,
	}
}
func (S *bestSeller) Instance(engName string) BestSellerFn {
	return map[string]BestSellerFn{
		"all":      S.GetBestSellerDefault,
		"health":   S.GetBestSellerHealth,
		"humanity": S.GetBestSellerHumanity,
		"business": S.GetBestSellerBusiness,
		"novel":    S.GetBestSellerNovel,
		"travel":   S.GetBestSellerTravel,
		"history":  S.GetBestSellerHistory,
		"selfdev":  S.GetBestSellerSelfDev,
		"dev":      S.GetBestSellerDev,
	}[engName]
}
func (S *bestSeller) GetBestSellerSelfDev() AladinResponse {
	return S.GetBestSeller("336")
}
func (S *bestSeller) GetBestSellerNovel() AladinResponse {
	return S.GetBestSeller("1")
}
func (S *bestSeller) GetBestSellerHumanity() AladinResponse {
	return S.GetBestSeller("656")
}
func (S *bestSeller) GetBestSellerHealth() AladinResponse {
	return S.GetBestSeller("55890")
}
func (S *bestSeller) GetBestSellerTravel() AladinResponse {
	return S.GetBestSeller("1196")
}
func (S *bestSeller) GetBestSellerDev() AladinResponse {
	return S.GetBestSeller("351")
}
func (S *bestSeller) GetBestSellerHistory() AladinResponse {
	return S.GetBestSeller("74")
}
func (S *bestSeller) GetBestSellerBusiness() AladinResponse {
	return S.GetBestSeller("170")
}
func (S *bestSeller) GetBestSellerDefault() AladinResponse {
	return S.GetBestSeller("0")
}
func (S *bestSeller) GetBestSeller(categoryId string) AladinResponse {
	// year := strconv.Itoa(time.Now().Year())
	// month := strconv.Itoa(int(time.Now().Month()))
	// week := strconv.Itoa(getWeekOfMonth(time.Now()))

	pageNum := []string{"1", "2", "3", "4"}

	var result []AladinItem
	for _, p := range pageNum {
		b, err := S.RequestBestSeller(p, categoryId, "0", "0", "0")
		if err != nil {
			log.Println(err)
		}
		resp := S.Parse(b)
		result = append(result, resp...)
	}

	return AladinResponse{
		CatValue: categoryId,
		CatName:  S.FindKorName(categoryId),
		Items:    result,
	}
}

func (S *bestSeller) RequestBestSeller(pageNum, categoryId, year, month, week string) (io.ReadCloser, error) {
	rawUrl, err := url.Parse(aladinItemUrl)
	if err != nil {
		return nil, err
	}
	rawQuery := rawUrl.Query()
	rawQuery.Add("TTBKey", S.ttbKey)
	rawQuery.Add("QueryType", "Bestseller")
	rawQuery.Add("Start", pageNum)
	rawQuery.Add("Cover", "MidBig")
	rawQuery.Add("MaxResults", "100")
	rawQuery.Add("CategoryId", categoryId)
	rawQuery.Add("SearchTarget", "Book")
	rawQuery.Add("Output", "JS")
	rawQuery.Add("Version", "20131101")
	if year != "0" {
		rawQuery.Add("Year", year)
		rawQuery.Add("Month", month)
		rawQuery.Add("Week", week)
	}

	rawUrl.RawQuery = rawQuery.Encode()

	log.Println(rawUrl.String())

	resp, err := http.Get(rawUrl.String())
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
func (S *bestSeller) Parse(body io.ReadCloser) []AladinItem {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil
	}
	var jsonData AladinItems
	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		return nil
	}
	return jsonData.Item
}

func (S *bestSeller) GetCategory() *[]BookCat {
	return cat
}
func (S *bestSeller) FindKorName(value string) string {
	for _, c := range *cat {
		if c.value == value {
			return c.korName
		}
	}
	return ""
}
func (S *bestSeller) GetCatName() []string {
	var result []string
	for _, c := range *cat {
		result = append(result, c.engName)
	}
	return result
}

// func getWeekOfMonth(t time.Time) int {
// 	// 이번 달의 첫 번째 날 가져오기
// 	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
// 	// 첫 번째 날의 요일 (0: 일요일, 1: 월요일, ...)
// 	firstDayWeekday := int(firstDayOfMonth.Weekday())
// 	// 오늘의 일(day)
// 	day := t.Day()

// 	// 이번 달 첫 번째 주의 남은 일수 계산
// 	remainingDays := 7 - firstDayWeekday
// 	if remainingDays <= 0 {
// 		remainingDays += 7
// 	}

// 	// 오늘이 첫 번째 주에 포함되는지 확인
// 	if day <= remainingDays {
// 		return 1
// 	}

// 	// 첫 번째 주 이후의 주 계산
// 	return (day-remainingDays-1)/7 + 2
// }
