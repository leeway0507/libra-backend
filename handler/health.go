package handler

import (
	"log"
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check: %s %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200"))
}
