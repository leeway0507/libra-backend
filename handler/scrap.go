package handler

import (
	"encoding/json"
	"libra-backend/pkg/scrap"
	"log"
	"net/http"
)

func HandleScrap(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	log.Printf("query: %#+v\n", query)
	isbn := r.PathValue("isbn")
	if isbn == "" {
		http.Error(w, "ISBN not provided", http.StatusBadRequest)
		return
	}

	engine := scrap.NewYangcheon(isbn)
	testData := scrap.NewLocalTest(engine)
	data := testData.ExtractDataFromLocal()

	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
