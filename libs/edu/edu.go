package edu

import (
	"fmt"
	"libra-backend/utils"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractData() {

	f, err := os.Open("./test.html")
	if err != nil {
		log.Println(err)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	bookCode := doc.Find("#position > tbody > tr > td:nth-child(3)").Text()
	bookStatus := doc.Find("#position > tbody > tr > td:nth-child(5)").Text()
	fmt.Println(bookCode)
	sp := strings.Split(strings.ReplaceAll(bookStatus, " ", ""), "\n")
	sp2 := utils.RemoveEmptyStringInSlice(sp)
	log.Println(strings.Join(sp2, " "))

}
