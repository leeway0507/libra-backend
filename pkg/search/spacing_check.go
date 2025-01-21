package search

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const PNU_URL = "http://speller.cs.pusan.ac.kr/results"

func GetSpacingCheck(keyword string) string {
	form := url.Values{}
	form.Add("text1", keyword)

	resp, err := http.Post(PNU_URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		// handle error
		log.Println(err)
		return keyword
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return keyword
	}
	return extractCorrectedWord(keyword, doc)
}

// 요청 결과 html으로 저장
func fetchSpacingCheckToHtml(path, keyword string) {
	htmlPath := filepath.Join(path, keyword+".html")

	// if html exists do nothing
	if _, err := os.Stat(htmlPath); err != nil {
		if os.IsExist(err) {
			log.Println(keyword, ":the html already exist")
			return
		}
	}

	form := url.Values{}
	form.Add("text1", keyword)

	resp, err := http.Post(PNU_URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	f, err := os.Create(htmlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

}

// 로컬에 저장한 html 불러오기
func readHtml(path string) *goquery.Document {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

type ErrorInfo struct {
	Help          string `json:"help"`
	ErrorIdx      int    `json:"errorIdx"`
	CorrectMethod int    `json:"correctMethod"`
	Start         int    `json:"start"`
	ErrMsg        string `json:"errMsg"`
	End           int    `json:"end"`
	OrgStr        string `json:"orgStr"`
	CandWord      string `json:"candWord"`
}

type GrammarSuggestion struct {
	Str     string      `json:"str"`
	ErrInfo []ErrorInfo `json:"errInfo"`
	Idx     int         `json:"idx"`
}

// html에서 GrammarSuggestion 추출
func extractResult(doc *goquery.Document) ([]GrammarSuggestion, error) {
	var grammarSugestions []GrammarSuggestion
	var errorS error
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if i == 2 {
			// 정규식을 사용하여 JSON 데이터 추출
			re := regexp.MustCompile(`data = (\[.*?\]);`)
			matches := re.FindStringSubmatch(s.Text())
			if len(matches) < 2 {
				errorS = fmt.Errorf("no matched data")
				return
			}

			jsonData := matches[1]
			err := json.Unmarshal([]byte(jsonData), &grammarSugestions)
			if err != nil {
				errorS = fmt.Errorf("error parsing JSON: %v", err)
				return
			}
		}
	})
	return grammarSugestions, errorS
}

// 수정된 키워드 반환
func extractCorrectedWord(keyword string, doc *goquery.Document) string {
	grammarSuggestion, err := extractResult(doc)
	correctedWord := keyword

	if err != nil {
		fmt.Println(err)
		return correctedWord
	}
	if len(grammarSuggestion) == 0 {
		return correctedWord
	}

	for _, erri := range grammarSuggestion[0].ErrInfo {
		if erri.Help == "띄어쓰기 오류입니다. 대치어를 참고하여 띄어 쓰도록 합니다." {
			candWords := strings.Split(erri.CandWord, "|")
			correctedWord = strings.Replace(correctedWord, erri.OrgStr, candWords[0], -1)
		}

	}
	return correctedWord
}
