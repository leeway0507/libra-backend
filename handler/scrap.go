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

	log.Printf("request info \n %#+v\n libCode:%s\n isbn:%s\n", engine.GetDistrict(), libCode, isbn)

	reader, err := engine.Request()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := engine.ExtractData(reader)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Not found info matched with ISBN", http.StatusNotFound)
		return
	}

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
