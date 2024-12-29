package handler

import (
	"encoding/json"
	"libra-backend/pkg/scrap"
	"log"
	"net/http"
)

func HandleScrap(w http.ResponseWriter, r *http.Request) {
	libCode := r.PathValue("libCode")
	isbn := r.PathValue("isbn")
	if isbn == "" {
		http.Error(w, "ISBN not provided", http.StatusBadRequest)
		return
	}

	engine := scrap.GetInstance(libCode)(isbn)
	reader, err := engine.Request()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := engine.ExtractData(reader)
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
