package lib_scrap

import (
	"libra-backend/model"
	"libra-backend/utils"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type education struct {
}

func NewEduction() *education {
	return &education{}
}

func (e *education) ExtractData() *model.LibBookStatus {

	f, err := os.Open("./test.html")
	if err != nil {
		log.Println(err)
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	bookCode := doc.Find("#position > tbody > tr > td:nth-child(3)").Text()
	bookStatusRaw := doc.Find("#position > tbody > tr > td:nth-child(5)").Text()
	var bookStatus string = bookStatusRaw
	if bookStatusRaw != "대출가능" {
		ss := strings.Split(strings.ReplaceAll(bookStatusRaw, " ", ""), "\n")
		ss2 := utils.RemoveEmptyStringInSlice(ss)
		bookStatus = strings.Join(ss2, " ")
	}

	return &model.LibBookStatus{
		LibType:    "education",
		Isbn:       "1234",
		LibName:    "서울특별시교육청양천도서관",
		BookCode:   bookCode,
		BookStatus: bookStatus,
	}

}
